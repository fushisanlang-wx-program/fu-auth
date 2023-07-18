package dao

import (
	"FuAuth/app/logger"
	"FuAuth/app/model"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

//

func VerifyOpenIdExist(baseStruct model.BaseStruct, openId string, ctx g.Ctx) bool {

	// if openId is in db,return true.else return false
	var (
		key = openId
	)
	sqlStr := "select count(1) as a from user where openid = ? ;"
	userExistCount, err := baseStruct.FuDb.GetOne(ctx, sqlStr, key)

	if err != nil {
		logger.LogError(baseStruct, gconv.String(err), ctx)
		return true
	} else {
		openIdExist := gconv.Int(userExistCount["a"])

		if openIdExist == 0 {
			return false
		} else if openIdExist == 1 {
			return true
		} else {
			logger.LogWarn(baseStruct, openId+" is more than 1,check it.", ctx)
			logger.LogInfo(baseStruct, openId+" is more than 1,check it.", ctx)
			return true
		}
	}

}

func VerifyUserExist(baseStruct model.BaseStruct, userName string, ctx g.Ctx) bool {

	// if username is in db,return true.else return false
	var (
		key = userName
	)
	sqlStr := "select count(1) as a from user where user = ? ;"
	userExistCount, err := baseStruct.FuDb.GetOne(ctx, sqlStr, key)

	if err != nil {
		logger.LogError(baseStruct, gconv.String(err), ctx)

		return true
	} else {
		userExist := gconv.Int(userExistCount["a"])

		if userExist == 0 {
			return false
		} else if userExist == 1 {
			return true
		} else {
			logger.LogWarn(baseStruct, userName+" is more than 1,check it.", ctx)
			logger.LogInfo(baseStruct, userName+" is more than 1,check it.", ctx)
			return true
		}
	}

}

func VerifyUser(baseStruct model.BaseStruct, UserName, UserPass string, ctx g.Ctx) bool {
	// VerifyUser name and pass.same return true,else return fasle
	sqlStr := "select count(1) as VerifyUserStatus from user where user = ? and pass = ? ;"
	verifyUserStatus, err := baseStruct.FuDb.GetOne(ctx, sqlStr, UserName, UserPass)
	if err != nil {
		logger.LogError(baseStruct, "VerifyUser err,UserName is "+gconv.String(UserName), ctx)

		return false
	} else {
		VerifyUserStatus := gconv.Int(verifyUserStatus["VerifyUserStatus"])
		if VerifyUserStatus == 0 {
			return false

		} else {
			return true
		}
	}

}

func UserAccessTimeRefresh(baseStruct model.BaseStruct, openId string, ctx g.Ctx) {
	_, err := baseStruct.FuDb.Update(ctx, "user", "lastaccesstime=CURRENT_TIMESTAMP", "openId=?", openId)
	logger.LogError(baseStruct, "Refresh user last Access Time  err,Uid is "+openId, ctx)
	if err != nil {
		logger.LogError(baseStruct, "Refresh user last Access Time  err,Uid is "+openId, ctx)

	}
}

func GetUserName(baseStruct model.BaseStruct, openId string, ctx g.Ctx) string {
	sqlStr := "select user from user where openid = ?"
	userInfo, err := baseStruct.FuDb.GetOne(ctx, sqlStr, openId)
	if err != nil {
		logger.LogError(baseStruct, "get user name err,openId is "+openId, ctx)
		return ""
	} else {
		UserName := gconv.String(userInfo["user"])
		return UserName
	}

}

func GetUserId(baseStruct model.BaseStruct, UserName string, ctx g.Ctx) int {
	sqlStr := "select id from user where user = ?"
	userInfo, err := baseStruct.FuDb.GetOne(ctx, sqlStr, UserName)
	if err != nil {
		logger.LogError(baseStruct, "get user info err,UserName is "+gconv.String(UserName), ctx)

		return 0
	} else {
		UserId := gconv.Int(userInfo["id"])
		return UserId
	}

}

func GetUserInfo(baseStruct model.BaseStruct, userOpenId string, ctx g.Ctx) model.UserInfoWithOpenId {

	sqlStr := "select * from user where openid = ?"
	userInfo, err := baseStruct.FuDb.GetOne(ctx, sqlStr, userOpenId)

	if err != nil {
		logger.LogError(baseStruct, "get user info err,user openid is "+userOpenId, ctx)

		UserInfo := model.UserInfoWithOpenId{}
		return UserInfo
	} else {
		UserInfo := model.UserInfoWithOpenId{
			UserOpenId: userOpenId,
			UserName:   gconv.String(userInfo["user"]),
		}
		return UserInfo
	}
}

func RegisterUser(baseStruct model.BaseStruct, UserInfo model.UserInfoWithOpenId, ctx g.Ctx) model.UserInfoWithOpenId {

	_, err := baseStruct.FuDb.Insert(ctx, "user", gdb.Map{
		"user":   UserInfo.UserName,
		"openid": UserInfo.UserOpenId,
	})
	if err != nil {
		logger.LogError(baseStruct, gconv.String(err), ctx)

		NilUser := model.UserInfoWithOpenId{}
		return NilUser
	} else {

		UserInfo = GetUserInfo(baseStruct, UserInfo.UserOpenId, ctx)
		return UserInfo
	}
}
