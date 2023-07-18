package api

import (
	"FuAuth/app/logger"
	"FuAuth/app/model"
	"FuAuth/app/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx = gctx.New()
)

func returnErrCode(r *ghttp.Request, code int, msg string) {
	r.Response.Status = code
	r.Response.WriteJson(g.Map{
		"message": msg,
	})
}

func GetOpenId(r *ghttp.Request, baseStruct model.BaseStruct) {

	CodeStr := r.Get("code").String()

	if CodeStr == "" {
		returnErrCode(r, 417, "数据空")
	} else {
		openId := service.GetOpenId(baseStruct, CodeStr, ctx)
		err := r.Session.Set("openId", openId.OpenId)
		if err != nil {
			logger.LogError(baseStruct, "set session openid as "+openId.OpenId+" err!", ctx)
		}
		if openId.Code == 200 {
			err := r.Session.Set("userName", openId.UserName)
			if err != nil {
				logger.LogError(baseStruct, "set session userName as "+openId.UserName+" err!", ctx)
			}
		}
		r.Response.WriteJson(openId)
	}
}

func UserRegister(r *ghttp.Request, baseStruct model.BaseStruct) {

	UserName := r.Get("UserName").String()
	UserCode := r.Get("code").String()
	// UserOpenId := r.Get("UserOpenId").String()
	UserOpenId := service.GetOpenId(baseStruct, UserCode, ctx).OpenId
	if UserName == "" || UserOpenId == "" {

		returnErrCode(r, 417, "用户注册失败，数据空")
	} else if service.VerifyOpenIdExist(baseStruct, UserOpenId, ctx) == true {
		returnErrCode(r, 423, "请勿重复注册")

	} else if service.VerifyUserExist(baseStruct, UserName, ctx) == true {
		returnErrCode(r, 423, "用户注册失败，用户名已存在")
	} else {
		UserInfo := service.RegisterUser(baseStruct, UserName, UserOpenId, ctx)
		err := r.Session.Set("userName", UserName)
		if err != nil {
			logger.LogError(baseStruct, "set session userName as "+UserName+" err!", ctx)
		}
		r.Response.WriteJson(UserInfo)

	}
}
