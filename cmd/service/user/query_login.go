package user

import (
	"ByteDance_5th/cmd/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util"
	"errors"
	"log"
)

const MaxNameLen = 100

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type QueryUserLoginFlow struct {
	username string
	password string

	data   *LoginResponse
	userid int64
	token  string
}

func QueryUserLogin(username, password string) (*LoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Operation()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

func (q *QueryUserLoginFlow) Operation() (*LoginResponse, error) {
	if err := q.checkName(); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		//log.Println("获取数据失败", err)
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

// 检查用户名与密码合法性
func (q *QueryUserLoginFlow) checkName() error {
	if q.username == "" {
		return errors.New(errortype.UserNameEmptyErr)
	}
	if len(q.username) > MaxNameLen {
		return errors.New(errortype.UserNameOverMaxLenErr)
	}
	if q.password == "" {
		errors.New(errortype.PasswordEmptyErr)
	}
	return nil
}

// 获取数据
func (q *QueryUserLoginFlow) prepareData() error {
	LoginDAO := models.NewLoginDao()
	var user models.User
	if err := LoginDAO.UserLogin(q.username, q.password, &user); err != nil {
		//log.Println("获取数据-A", err)
		return err
	}
	q.userid = user.UserInfoId

	//获取token
	var (
		token string
		err   error
	)
	if token, err = util.GenerateToken(user); err != nil {
		//log.Println("获取数据-B", err)
		return err
	}
	q.token = token
	return nil
}

// 打包
func (q *QueryUserLoginFlow) packData() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
