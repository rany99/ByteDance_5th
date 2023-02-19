package login

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util"
	"errors"
)

// PostUserLogin 注册用户并得到token和id
func PostUserLogin(username, password string) (*LoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Operation()
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

func (q *PostUserLoginFlow) Operation() (*LoginResponse, error) {
	//对参数进行合法性验证
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}

	//更新数据到数据库
	if err := q.GetData(); err != nil {
		return nil, err
	}

	//打包response
	if err := q.PackResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *PostUserLoginFlow) CheckJSON() error {
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

func (q *PostUserLoginFlow) GetData() error {

	//准备好userInfo,默认name为username

	userLogin := models.User{
		Username: q.username,
		Password: q.password,
	}
	userInfo := models.UserInfo{
		Name:            q.username,
		Avatar:          util.GetAvatarUrl(),
		BackgroundImage: util.GetBackGroundUrl(),
		Signature:       util.GetSignature(),
		User:            &userLogin,
	}

	//判断用户名是否已经存在
	userLoginDAO := models.NewLoginDao()
	if userLoginDAO.UserAlreadyExist(q.username) {
		return errors.New(errortype.UserNameExistErr)
	}

	userInfoDAO := models.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userInfo)
	if err != nil {
		return err
	}

	//颁发token
	var token string
	if token, err = util.GenerateToken(userLogin); err != nil {
		return err
	}
	q.token = token
	q.userid = userInfo.Id
	return nil
}

func (q *PostUserLoginFlow) PackResponse() error {
	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
