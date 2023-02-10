package util

import (
	"ByteDance_5th/cache"
	"ByteDance_5th/config"
	"ByteDance_5th/models"
	"errors"
	"fmt"
	"log"
	"time"
)

// GetVideoUrl 返回url
func GetVideoUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/videos/%s", config.Conf.SE.IP, config.Conf.SE.Port, fileName)
}

// GetImageUrl 返回url
func GetImageUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/images/%s", config.Conf.SE.IP, config.Conf.SE.Port, fileName)
}

// NewUnicFileName 生成文件名
func NewUnicFileName(userid int64) string {
	var count int64
	if err := models.NewVideoDao().QueryVideoCntByUserId(userid, &count); err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userid, count)
}

// FillVideos 更新视频作者信息
func FillVideos(userid int64, videos *[]*models.Video) (*time.Time, error) {
	videosLen := len(*videos)
	if videos == nil || videosLen == 0 {
		return nil, errors.New("FillVideos：传入videos列表为空")
	}
	dao := models.NewUserInfoDAO()
	p := cache.NewProxyIndexMap()
	latestTime := (*videos)[videosLen-1].CreatedAt
	for i := 0; i < videosLen; i++ {
		var author models.UserInfo
		if err := dao.QueryUserInfoById((*videos)[i].UserInfoId, &author); err != nil {
			continue
		}
		author.IsFollow = p.GetAFollowB(userid, author.Id)
		(*videos)[i].Author = author
		if userid > 0 {
			(*videos)[i].IsFavorite = p.GetVideoFavor(userid, (*videos)[i].Id)
		}
	}
	return &latestTime, nil
}

// SaveImageFromVideo 截取视频封面
//func SaveImageFromVideo(name string, isDebug bool) error {
//	v2i := NewVideoToCover()
//	if isDebug {
//		v2i.Debug()
//	}
//	v2i.InputPath = filepath.Join(config.Conf.Path.StaticSourcePath, name+GetDefaultVideoSuffix())
//	v2i.OutputPath = filepath.Join(config.Conf.Path.StaticSourcePath, name+GetDefaultImageSuffix())
//	v2i.FrameCount = 1
//	queryString, err := v2i.GetQueryString()
//	if err != nil {
//		return err
//	}
//	return v2i.ExecCommand(queryString)
//}
