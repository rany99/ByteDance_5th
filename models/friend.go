package models

const AvatarCnt = 16

type Friend struct {
	UserInfo UserInfo
	Avatar   string `json:"avatar"`
}
