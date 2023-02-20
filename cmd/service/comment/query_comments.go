package comment

import (
	models2 "ByteDance_5th/cmd/models"
	"ByteDance_5th/pkg/errortype"
	"errors"
)

type CommentsResponse struct {
	Comments []*models2.Comment `json:"comment_list"`
}

type QueryCommentListFlow struct {
	//uid         int64
	vid         int64
	comments    []*models2.Comment
	commentList *CommentsResponse
}

func QueryCommentList(vid int64) (*CommentsResponse, error) {
	return NewQueryCommentListFlow(vid).Operation()
}

func NewQueryCommentListFlow(vid int64) *QueryCommentListFlow {
	return &QueryCommentListFlow{
		//uid: uid,
		vid: vid,
	}
}

func (q *QueryCommentListFlow) Operation() (*CommentsResponse, error) {
	if err := q.CheckJson(); err != nil {
		//log.Println("QueryCommentListFlow:checkJSON失败")
		return nil, err
	}
	if err := q.GetData(); err != nil {
		//log.Println("QueryCommentListFlow:getData失败")
		return nil, err
	}
	if err := q.PackData(); err != nil {
		//log.Println("QueryCommentListFlow:PackData失败")
		return nil, err
	}
	return q.commentList, nil
}

func (q *QueryCommentListFlow) CheckJson() error {
	//判断用户是否存在
	//if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
	//	return err
	//}
	//判断视频是否存在
	if !models2.NewVideoDao().VideoAlreadyExist(q.vid) {
		return errors.New(errortype.VideoNoExistErr)
	}
	return nil
}

func (q *QueryCommentListFlow) GetData() error {
	if err := models2.NewCommentDao().QueryCommentListByVideoId(q.vid, &q.comments); err != nil {
		return err
	}
	if err := FillCommentList(&q.comments); err != nil {
		return errors.New(errortype.VideoHasNoCommentErr)
	}
	return nil
}

func (q *QueryCommentListFlow) PackData() error {
	q.commentList = &CommentsResponse{Comments: q.comments}
	return nil
}

func FillCommentList(comments *[]*models2.Comment) error {
	if comments == nil {
		return errors.New("FillCommentListFields" + errortype.PointerIsNilErr)
	}
	commentsLen := len(*comments)
	if commentsLen == 0 {
		return errors.New(errortype.VideoListEmptyErr)
	}
	userInfoDAO := models2.NewUserInfoDAO()
	for _, c := range *comments {
		_ = userInfoDAO.QueryUserInfoById(c.UserInfoId, &c.User)
		c.CreateDate = c.CreatedAt.Format("1-2")
	}
	return nil
}
