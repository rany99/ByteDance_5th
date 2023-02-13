package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
)

type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	Name          string      `json:"name" gorm:"name,omitempty"`
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *User       `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                     //用户与评论的一对多
}

type UserInfoDao struct {
}

var (
	userInfoDao  *UserInfoDao
	userInfoOnce sync.Once
)

// NewUserInfoDAO 创建DAO
func NewUserInfoDAO() *UserInfoDao {
	userInfoOnce.Do(func() {
		userInfoDao = new(UserInfoDao)
	})
	return userInfoDao
}

// QueryUserInfoById 查询用户
func (u *UserInfoDao) QueryUserInfoById(id int64, userinfo *UserInfo) error {
	if userinfo == nil {
		return errors.New("传入UserInfo指针为空")
	}
	DB.Where("id = ?", id).Select([]string{"id", "name", "follow_count", "is_follow"}).First(userinfo)
	if userinfo.Id == 0 {
		return errors.New("查询不到该用户")
	}
	return nil
}

// AddUserInfo 将UserInfo指针信息写入数据库
func (u *UserInfoDao) AddUserInfo(userinfo *UserInfo) error {
	if userinfo == nil {
		return errors.New("传入userinfo指针为空")
	}
	return DB.Create(userinfo).Error
}

// NoAFollowB 取消关注关系
func (u *UserInfoDao) NoAFollowB(a, b int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", a).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", b).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_relations` WHERE user_info_id=? AND follow_id=?", a, b).Error; err != nil {
			return err
		}
		return nil
	})
}

// AFollowB 建立关注关系
func (u *UserInfoDao) AFollowB(a, b int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count+1 WHERE id = ?", a).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count+1 WHERE id = ?", b).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `user_relations` (`user_info_id`,`follow_id`) VALUES (?,?)", a, b).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetFollowsById 获取关注列表
func (u *UserInfoDao) GetFollowsById(id int64, userList *[]*UserInfo) error {
	if userList == nil {
		return errors.New("传入指针为空")
	}
	if err := DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.user_info_id = ? AND r.follow_id = u.id", id).Scan(userList).Error; err != nil {
		return err
	}
	log.Println("GetFollowListById", len(*userList))
	return nil
}

// GetFansById 获取粉丝列表
func (u *UserInfoDao) GetFansById(id int64, userList *[]*UserInfo) error {
	if userList == nil {
		return errors.New("传入指针为空")
	}
	if err := DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.follow_id = ? AND r.user_info_id = u.id", id).Scan(userList).Error; err != nil {
		return err
	}
	return nil
}

// IsUserInfoExist 判断用户表中是否存在该id的用户
func (u *UserInfoDao) IsUserInfoExist(id int64) error {
	var userInfo UserInfo
	if err := DB.Where("id = ?", id).Select("id").First(&userInfo).Error; err != nil {
		return err
	}
	if userInfo.Id == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
