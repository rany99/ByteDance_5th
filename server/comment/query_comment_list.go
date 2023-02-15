package comment

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"errors"
)

type CList struct {
	Comments []*models.Comment `json:"comment_list"`
}

type QueryCommentListFlow struct {
	uid         int64
	vid         int64
	comments    []*models.Comment
	commentList *CList
}

func QueryCommentList(uid, vid int64) (*CList, error) {
	return NewQueryCommentListFlow(uid, vid).Operation()
}

func NewQueryCommentListFlow(uid int64, vid int64) *QueryCommentListFlow {
	return &QueryCommentListFlow{
		uid: uid,
		vid: vid,
	}
}

func (q *QueryCommentListFlow) Operation() (*CList, error) {
	if err := q.checkJson(); err != nil {
		//log.Println("QueryCommentListFlow:checkJSON失败")
		return nil, err
	}
	if err := q.getData(); err != nil {
		//log.Println("QueryCommentListFlow:getData失败")
		return nil, err
	}
	if err := q.packData(); err != nil {
		//log.Println("QueryCommentListFlow:PackData失败")
		return nil, err
	}
	return q.commentList, nil
}

func (q *QueryCommentListFlow) checkJson() error {
	//判断用户是否存在
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	//判断视频是否存在
	if !models.NewVideoDao().VideoAlreadyExist(q.vid) {
		return errors.New(errortype.VideoNoExistErr)
	}
	return nil
}

func (q *QueryCommentListFlow) getData() error {
	if err := models.NewCommentDao().QueryCommentListByVideoId(q.vid, &q.comments); err != nil {
		return err
	}
	if err := FillCommentList(&q.comments); err != nil {
		return errors.New(errortype.VideoHasNoCommentErr)
	}
	return nil
}

func (q *QueryCommentListFlow) packData() error {
	q.commentList = &CList{Comments: q.comments}
	return nil
}

func FillCommentList(comments *[]*models.Comment) error {
	if comments == nil {
		return errors.New("FillCommentListFields" + errortype.PointerIsNilErr)
	}
	commentsLen := len(*comments)
	if commentsLen == 0 {
		return errors.New(errortype.VideoListEmptyErr)
	}
	userInfoDAO := models.NewUserInfoDAO()
	for _, c := range *comments {
		_ = userInfoDAO.QueryUserInfoById(c.UserInfoId, &c.User)
		c.CreateDate = c.CreatedAt.Format("1-2")
	}
	return nil
}
