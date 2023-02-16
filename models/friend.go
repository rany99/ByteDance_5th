package models

const AvatarCnt = 16

type Friend struct {
	UserInfo
	Avatar string `json:"avatar"`
}
