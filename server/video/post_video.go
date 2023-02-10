package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util"
	"fmt"
)

type PostVideoFlow struct {
	VideoName string
	CoverName string
	Title     string
	UserId    int64
	Video     *models.Video
}

func (p *PostVideoFlow) Do() error {

	p.GenerateUrl()
	if err := p.publish(); err != nil {
		return err
	}
	return nil
}

func (p *PostVideoFlow) GenerateUrl() {
	p.VideoName = util.GetVideoUrl(p.VideoName)
	p.CoverName = util.GetImageUrl(p.CoverName)
}

func (p *PostVideoFlow) publish() error {
	video := &models.Video{
		UserInfoId: p.UserId,
		PlayUrl:    p.VideoName,
		CoverUrl:   p.CoverName,
		Title:      p.Title,
	}
	fmt.Println("publish-UserInfoId", video.UserInfoId)
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
	return NewPostVideoFlow(userId, videoName, coverName, title).Do()
}
