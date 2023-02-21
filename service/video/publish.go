package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/config"
	"fmt"
)

type PostVideoFlow struct {
	VideoName string
	CoverName string
	Title     string
	UserId    int64
	Video     *models.Video
}

func (p *PostVideoFlow) Opeartion() error {
	p.GenerateUrl()
	if err := p.publish(); err != nil {
		return err
	}
	return nil
}

func (p *PostVideoFlow) GenerateUrl() {
	//log.Println("VideoName", p.VideoName)
	//log.Println("CoverName", p.CoverName)
	p.VideoName = GetVideoUrl(p.VideoName)
	p.CoverName = GetImageUrl(p.CoverName)
}

func (p *PostVideoFlow) publish() error {
	video := &models.Video{
		UserInfoId: p.UserId,
		PlayUrl:    p.VideoName,
		CoverUrl:   p.CoverName,
		Title:      p.Title,
	}
	//log.Println("publish-UserInfoId", video.UserInfoId)
	return models.NewVideoDao().AddVideoToDB(video)
}

func NewPostVideoFlow(userId int64, videoName, coverName, title string) *PostVideoFlow {
	return &PostVideoFlow{
		VideoName: videoName,
		CoverName: coverName,
		UserId:    userId,
		Title:     title,
	}
}

// PostVideo 发布视频
func PostVideo(userId int64, videoName, coverName, title string) error {
	return NewPostVideoFlow(userId, videoName, coverName, title).Opeartion()
}

// GetVideoUrl 返回url
func GetVideoUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/static/video/%s", config.Conf.SE.IP, config.Conf.SE.Port, fileName)
}

// GetImageUrl 返回url
func GetImageUrl(fileName string) string {
	return fmt.Sprintf("http://%s:%d/static/cover/%s", config.Conf.SE.IP, config.Conf.SE.Port, fileName)
}
