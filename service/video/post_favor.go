package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util/cache"
	"errors"
)

type PostFavorFlow struct {
	uid        int64
	vid        int64
	actionType int64
}

func PostFavor(uid, vid, actionType int64) error {
	return NewPostFavorFlow(uid, vid, actionType).Do()
}

func NewPostFavorFlow(uid, vid, actionType int64) *PostFavorFlow {
	return &PostFavorFlow{
		uid:        uid,
		vid:        vid,
		actionType: actionType,
	}
}

func (p *PostFavorFlow) Do() error {
	if err := p.IsUserExist(); err != nil {
		return err
	}
	if err := p.IsActionTypeLegal(); err != nil {
		return err
	}
	if p.actionType == 1 {
		if err := p.AddOne(); err != nil {
			return err
		}
	} else {
		if err := p.SubOne(); err != nil {
			return err
		}
	}
	return nil
}

// IsUserExist 判断用户是否存在
func (p *PostFavorFlow) IsUserExist() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(p.uid); err != nil {
		return err
	}
	return nil
}

// IsActionTypeLegal 判断actionType时候合法
func (p *PostFavorFlow) IsActionTypeLegal() error {
	if p.actionType == 1 || p.actionType == 2 {
		return nil
	}
	return errors.New(errortype.PostFavorActionTypeErr)
}

// AddOne 执行点赞操作
func (p *PostFavorFlow) AddOne() error {
	if err := models.NewVideoDao().FavoriteCountAddOneByVideoId(p.uid, p.vid); err != nil {
		return errors.New(errortype.AlreadyPostFavorErr)
	}
	cache.NewProxyIndexMap().SetVideoFavor(p.uid, p.vid, true)
	return nil
}

// SubOne 执行取消点赞操作
func (p *PostFavorFlow) SubOne() error {
	if err := models.NewVideoDao().FavoriteCountSubOneByVideoId(p.uid, p.vid); err != nil {
		if err != nil {
			errors.New(errortype.FavorCountZeroErr)
		}
	}
	cache.NewProxyIndexMap().SetVideoFavor(p.uid, p.vid, false)
	return nil
}
