package login

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util"
	"errors"
)

// PostUserLogin 注册用户并得到token和id
func PostUserLogin(username, password string) (*LoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Do()
}

func NewPostUserLoginFlow(username, password string) *PostUserLoginFlow {
	return &PostUserLoginFlow{username: username, password: password}
}

type PostUserLoginFlow struct {
	username string
	password string

	data   *LoginResponse
	userid int64
	token  string
}

func (q *PostUserLoginFlow) Do() (*LoginResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}

	//更新数据到数据库
	if err := q.updateData(); err != nil {
		return nil, err
	}

	//打包response
	if err := q.packResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *PostUserLoginFlow) checkNum() error {
	if q.username == "" {
		return errors.New(errortype.UserNameEmptyErr)
	}
	if len(q.username) > MaxNameLen {
		return errors.New(errortype.UserNameOverMaxLenErr)
	}
	if q.password == "" {
		return errors.New(errortype.PasswordEmptyErr)
	}
	return nil
}

func (q *PostUserLoginFlow) updateData() error {

	//准备好userInfo,默认name为username

	userLogin := models.User{Username: q.username, Password: q.password}
	userinfo := models.UserInfo{User: &userLogin, Name: q.username}

	//判断用户名是否已经存在
	userLoginDAO := models.NewLoginDao()
	if userLoginDAO.UserAlreadyExist(q.username) {
		return errors.New(errortype.UserNameExistErr)
	}

	userInfoDAO := models.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return err
	}

	//颁发token
	var token string
	if token, err = util.GenerateToken(userLogin); err != nil {
		return err
	}
	q.token = token
	q.userid = userinfo.Id
	return nil
}

func (q *PostUserLoginFlow) packResponse() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
