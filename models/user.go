package models

import (
	"errors"
	"log"
	"sync"
)

type LoginDao struct {
}

// User 登陆信息, UserInfo 获取用户follow等详细信息的外键
type User struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

var (
	userLoginDao *LoginDao
	//避免重复注册
	LoginOnce sync.Once
)

// 登陆
func (u *LoginDao) UserLogin(username, password string, login *User) error {
	if login == nil {
		return errors.New("输入User结构体为空")
	}
	DB.Where("username = ?", username).First(login)
	if login.Id == 0 {
		return errors.New("账户输入错误或用户不存在")
	}
	DB.Where("username = ? AND password = ?", username, password).First(login)
	if login.Id == 0 {
		return errors.New("密码输入错误")
	}
	return nil
}

func NewLoginDao() *LoginDao {
	LoginOnce.Do(func() {
		userLoginDao = new(LoginDao)
	})
	return userLoginDao
}

// 查询账户是否已经存在
func (u *LoginDao) UserAlreadyExist(username string) bool {
	var user User
	log.Println(username)
	DB.Where("username = ?", username).First(&user)
	if user.Id == 0 {
		log.Println("无用户记录")
		return false
	}
	return true
}
