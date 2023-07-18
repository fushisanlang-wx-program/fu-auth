package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"FuAuth/app/dao"
	"FuAuth/app/logger"
	"FuAuth/app/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/gogf/gf/v2/util/gconv"
)

func VerifyUserExist(baseStruct model.BaseStruct, userName string, ctx g.Ctx) bool {

	return dao.VerifyUserExist(baseStruct, userName, ctx)
}

func VerifyOpenIdExist(baseStruct model.BaseStruct, OpenId string, ctx g.Ctx) bool {

	return dao.VerifyOpenIdExist(baseStruct, OpenId, ctx)
}
func VerifyUser(baseStruct model.BaseStruct, UserName, UserPass string, ctx g.Ctx) bool {
	return dao.VerifyUser(baseStruct, UserName, UserPass, ctx)
}

func UserAccessTimeRefresh(baseStruct model.BaseStruct, openId string, ctx g.Ctx) {
	dao.UserAccessTimeRefresh(baseStruct, openId, ctx)
}

func GetUserId(baseStruct model.BaseStruct, UserName string, ctx g.Ctx) int {

	return dao.GetUserId(baseStruct, UserName, ctx)
}
func GetUserName(baseStruct model.BaseStruct, OpenId string, ctx g.Ctx) string {

	return dao.GetUserName(baseStruct, OpenId, ctx)
}
func RegisterUser(baseStruct model.BaseStruct, userName, userOpenId string, ctx g.Ctx) model.UserInfoWithOpenId {
	UserInfo := model.UserInfoWithOpenId{
		//UserId:             userId,
		UserName:   userName,
		UserOpenId: userOpenId,
		UserAdmin:  false,
		//UserCreateTime:     "",
		//UserLastAccessTime: "",
	}

	return dao.RegisterUser(baseStruct, UserInfo, ctx)
}
func GetOpenId(baseStruct model.BaseStruct, CodeStr string, ctx g.Ctx) model.OpenId {
	logger.LogInfo(baseStruct, "get codestr as "+CodeStr, ctx)
	GetStatus, OpenId := getOpenId(baseStruct, CodeStr, ctx)
	var code int
	var message, userName string

	if GetStatus {

		openIdExist := VerifyOpenIdExist(baseStruct, OpenId, ctx)

		if openIdExist {
			UserAccessTimeRefresh(baseStruct, OpenId, ctx)
			userName = GetUserName(baseStruct, OpenId, ctx)
			//用户存在
			code = 200
			message = ""
		} else {
			//用户不存在
			code = 401
			message = "用户未注册"
			userName = ""
		}

	} else {
		code = 502
		message = "微信服务器异常，请重试或联系开发者"
		userName = ""
	}
	var UserOpenId = model.OpenId{
		Code:     code,
		OpenId:   OpenId,
		Message:  message,
		UserName: userName,
	}
	logger.LogInfo(baseStruct, "get user openid as "+OpenId, ctx)
	return UserOpenId
}

func getOpenId(baseStruct model.BaseStruct, codeStr string, ctx g.Ctx) (bool, string) {
	jscode2sessionUrlPath, err := baseStruct.FuConf.Get(ctx, "wx.jscode2sessionUrlPath")
	if err != nil {
		logger.LogError(baseStruct, "urlPath not config.", ctx)

	}
	appid, err := baseStruct.FuConf.Get(ctx, "wx.appid")
	if err != nil {
		logger.LogError(baseStruct, "appid not config.", ctx)

	}
	secret, err := baseStruct.FuConf.Get(ctx, "wx.secret")
	if err != nil {
		logger.LogError(baseStruct, "secret not config.", ctx)

	}
	grantType, err := baseStruct.FuConf.Get(ctx, "wx.grant_type")
	if err != nil {
		logger.LogError(baseStruct, "grant_type not config.", ctx)

	}
	urlPath := gconv.String(jscode2sessionUrlPath)
	postString := "appid=" + gconv.String(appid) + "&secret=" + gconv.String(secret) + "&js_code=" + gconv.String(codeStr) + "&grant_type=" + gconv.String(grantType)
	postStringByte := []byte(postString)
	req, _ := http.NewRequest("Get", urlPath, bytes.NewBuffer(postStringByte))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, ""

	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		OpenId := transformation(resp)["openid"]
		return true, gconv.String(OpenId)
	} else {
		return false, ""
	}

}

func VerifyCode(baseStruct model.BaseStruct, r *ghttp.Request, ctx g.Ctx) model.OpenId {
	CodeStr := r.Get("code").String()

	return GetOpenId(baseStruct, CodeStr, ctx)
}

func VerifySession(baseStruct model.BaseStruct, r *ghttp.Request, ctx g.Ctx) (bool, string, string) {

	sessionData, err := r.Session.Data()
	if err != nil {
		return false, "", ""
	}
	var userStruct *model.UserInfoWithOpenId

	if gconv.Struct(sessionData, &userStruct) != nil {
		return false, "", ""
	}
	UserName := userStruct.UserName

	UserOpenId := userStruct.UserOpenId
	userNameTrue := GetUserName(baseStruct, UserOpenId, ctx)
	if UserOpenId == "" {
		return false, "", ""
	} else if UserName == userNameTrue {
		return true, UserOpenId, UserName
	} else {
		return false, "", ""
	}

}
func transformation(response *http.Response) map[string]interface{} {
	var result map[string]interface{}
	body, err := io.ReadAll(response.Body)
	if err == nil {
		json.Unmarshal([]byte(string(body)), &result)
	}

	return result
}
