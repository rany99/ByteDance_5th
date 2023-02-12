package comment

import (
	"ByteDance_5th/models"
	"errors"
	"log"
)

type CResponse struct {
	CRComment *models.Comment `json:"comment"`
}

type PostCommentFlow struct {
	uid         int64
	vid         int64
	commentId   int64
	actionType  int64
	commentText string
	comment     *models.Comment
	*CResponse
}

// PostComment 发布评论
func PostComment(uid, vid, commentId, actionType int64, commentText string) (*CResponse, error) {
	return NewPostCommentFlow(uid, vid, commentId, actionType, commentText).Do()
}

func NewPostCommentFlow(uid int64, vid int64, commentId int64, actionType int64, commentText string) *PostCommentFlow {
	return &PostCommentFlow{uid: uid, vid: vid, commentId: commentId, actionType: actionType, commentText: commentText}
}

func (p *PostCommentFlow) Do() (*CResponse, error) {
	if err := p.CheckJson(); err != nil {
		return nil, err
	}
	if err := p.GetData(); err != nil {
		return nil, err
	}
	log.Println("Do GetData:", p.comment.Content)
	if err := p.PackData(); err != nil {
		return nil, err
	}
	log.Println("Do PackData:", p.comment.Content)
	log.Println(p.CResponse.CRComment.Content)
	return p.CResponse, nil
}

// CheckJson 检查Json传入数据是否正确
func (p *PostCommentFlow) CheckJson() error {
	if err := p.CheckUid(); err != nil {
		return err
	}
	if err := p.CheckVid(); err != nil {
		return err
	}
	if err := p.CheckActionType(); err != nil {
		return err
	}
	return nil
}

// CheckUid 检查用户是否存在
func (p *PostCommentFlow) CheckUid() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(p.uid); err != nil {
		return err
	}
	return nil
}

// CheckVid 检查视频是否窜在
func (p *PostCommentFlow) CheckVid() error {
	if ok := models.NewVideoDao().VideoAlreadyExist(p.vid); !ok {
		return errors.New("当前videoId的视频不存在")
	}
	return nil
}

// CheckActionType 检查ActionType是否合法
func (p *PostCommentFlow) CheckActionType() error {
	if p.actionType == 1 || p.actionType == 2 {
		return nil
	}
	return errors.New("只能进行评论1或者删除2操作")
}

// GetData 获取数据
func (p *PostCommentFlow) GetData() error {
	var err error
	switch p.actionType {
	case 1: //创建
		p.comment, err = p.CreateComment()
	case 2: //删除
		p.comment, err = p.DeleteComment()
	default:
		return errors.New("只能进行评论1或者删除2操作")
	}
	log.Println("GetData:", p.comment.Content)
	return err
}

// CreateComment 创建评论
func (p *PostCommentFlow) CreateComment() (*models.Comment, error) {
	comment := models.Comment{
		UserInfoId: p.uid,
		VideoId:    p.vid,
		Content:    p.commentText,
	}
	log.Println("CreateComment")
	if err := models.NewCommentDao().CreateAndCntAddOne(&comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment 删除评论
func (p *PostCommentFlow) DeleteComment() (*models.Comment, error) {
	var comment models.Comment
	//确认评论是否存在，并获得待删除评论指针
	if err := models.NewCommentDao().QueryCommentById(p.commentId, &comment); err != nil {
		return nil, err
	}
	//删除评论
	if err := models.NewCommentDao().DeleteAndCntSubOne(p.commentId, p.vid); err != nil {
		return nil, err
	}
	return &comment, nil
}

// PackData 封装数据
func (p *PostCommentFlow) PackData() error {
	var userInfo models.UserInfo
	_ = models.NewUserInfoDAO().QueryUserInfoById(p.comment.UserInfoId, &userInfo)
	p.comment.User = userInfo
	_ = FillCommentFields(p.comment)
	p.CResponse = &CResponse{CRComment: p.comment}
	return nil
}

func FillCommentFields(comment *models.Comment) error {
	if comment == nil {
		return errors.New("FillCommentFields:传入Comment指针为空")
	}
	comment.CreateDate = comment.CreatedAt.Format("1-2")
	return nil
}
