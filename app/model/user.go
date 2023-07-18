package model

type UserInfo struct {
	UserId             int
	UserName           string
	UserPass           string
	UserAdmin          bool
	UserCreateTime     string
	UserLastAccessTime string
}

type OpenId struct {
	Code     int
	OpenId   string
	Message  string
	UserName string
}

type UserInfoWithOpenId struct {
	UserOpenId         string
	UserName           string
	UserAdmin          bool
	UserCreateTime     string
	UserLastAccessTime string
}
