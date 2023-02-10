package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

type VideoDao struct {
}

var (
	videoDao  *VideoDao
	videoOnce sync.Once
)

func NewVideoDao() *VideoDao {
	videoOnce.Do(func() {
		videoDao = new(VideoDao)
	})
	return videoDao
}

// AddVideoToDB 在数据库中添加新的视频
func (v *VideoDao) AddVideoToDB(video *Video) error {
	if video == nil {
		return errors.New("AddVideoToDB：传入Video指针为空")
	}
	return DB.Create(video).Error
}

// QueryVideoByVideoId 通过视频ID返回视频结构体
func (v *VideoDao) QueryVideoByVideoId(id int64, video *Video) error {
	if video == nil {
		return errors.New("QueryVideoById：传入Video指针为空")
	}
	return DB.Where("id = ?", id).Select([]string{
		"id",
		"user_info_id",
		"play_url",
		"cover_url",
		"favorite_count",
		"comment_count",
		"is_favorite",
		"title",
	}).First(video).Error
}

// QueryVideoCntByUserId 返回作者发布的视频数量
func (v *VideoDao) QueryVideoCntByUserId(id int64, cnt *int64) error {
	if cnt == nil {
		return errors.New("QueryVideoCntByUserId：cnt指针为空")
	}
	return DB.Model(&Video{}).Where("user_info_id = ?", id).Count(cnt).Error
}

// QueryVideoListByUserId 通过userid返回视频列表
func (v *VideoDao) QueryVideoListByUserId(id int64, list *[]*Video) error {
	if list == nil {
		return errors.New("QueryVideoListByUserId：list指针为空")
	}
	return DB.Where("user_info_id = ?", id).Select([]string{
		"id",
		"user_info_id",
		"play_url",
		"cover_url",
		"favorite_count",
		"comment_count",
		"is_favorite",
		"title",
	}).Find(list).Error
}

// QueryVideoListByLastTimeAndLimit 根据latestTime返回之前的Limit个视频
func (v *VideoDao) QueryVideoListByLastTimeAndLimit(latestTime time.Time, limit int, list *[]*Video) error {
	if list == nil {
		log.Println("QueryVideoListByLastTimeAndLimit：list指针为空")
		return errors.New("QueryVideoListByLastTimeAndLimit：list指针为空")
	}
	log.Println("latestTime:", latestTime)
	err := DB.Model(&Video{}).Where("created_at<?", latestTime).
		Order("created_at ASC").Limit(limit).
		Select([]string{"id", "user_info_id", "play_url", "cover_url", "favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).
		Find(list).Error
	log.Println("从数据库中获得的List长度为：", len(*list))
	return err
}

// FavoriteCountAddOneByVideoId 根据视频ID和用户ID将视频点赞数加一，并添加到user_favor_videos
func (v *VideoDao) FavoriteCountAddOneByVideoId(userid int64, videoid int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE videos SET favorite_count = favorite_count + 1 WHERE id = ?", videoid).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `user_favor_videos` (`user_info_id`,`video_id`) VALUES (?,?)", userid, videoid).Error; err != nil {
			return err
		}
		return nil
	})
}

// FavoriteCountSubOneByVideoId 根据视频ID和用户ID将视频点赞数减一，并从user_favor_videos删除
func (v *VideoDao) FavoriteCountSubOneByVideoId(userid int64, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE videos SET favorite_count = favorite_count + 1 WHERE id = ? AND favorite_count > 0", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_favor_videos` (`user_info_id`,`video_id`) VALUES (?,?)", userid, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryFavorListByUserId 获取用户点赞视频列表
func (v *VideoDao) QueryFavorListByUserId(userid, videoid int64, list *[]*Video) error {
	if err := DB.Raw("SELECT v.* FROM user_favor_videos u , videos v WHERE u.user_info_id = ? AND u.video_id = v.id", userid).Scan(list).Error; err != nil {
		return err
	}
	if CheckList(list) {
		return errors.New("用户点赞列表为空")
	}
	return nil
}

// CheckList 检测点赞视频列表是否为空
func CheckList(list *[]*Video) bool {
	return len(*list) == 0 || (*list)[0].Id == 0
}

// VideoAlreadyExist 检测视频是否存在
func (v *VideoDao) VideoAlreadyExist(id int64) bool {
	var video Video
	if err := DB.Where("id = ?", id).Select("id").First(&video).Error; err != nil {
		log.Println(err)
	}
	return video.Id != 0
}
