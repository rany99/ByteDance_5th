package comment

import (
	"ByteDance_5th/models"
	"errors"
	"log"
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
	return NewQueryCommentListFlow(uid, vid).Do()
}

func NewQueryCommentListFlow(uid int64, vid int64) *QueryCommentListFlow {
	return &QueryCommentListFlow{
		uid: uid,
		vid: vid,
	}
}

func (q *QueryCommentListFlow) Do() (*CList, error) {
	if err := q.checkJson(); err != nil {
		log.Println("QueryCommentListFlow:checkJSON失败")
		return nil, err
	}
	if err := q.getData(); err != nil {
		log.Println("QueryCommentListFlow:getData失败")
		return nil, err
	}
	if err := q.packData(); err != nil {
		log.Println("QueryCommentListFlow:PackData失败")
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
		return errors.New("视频不存在")
	}
	return nil
}

func (q *QueryCommentListFlow) getData() error {
	if err := models.NewCommentDao().QueryCommentListByVideoId(q.vid, &q.comments); err != nil {
		return err
	}
	if err := FillCommentListFields(&q.comments); err != nil {
		return errors.New("还没有人发现这里，赶紧抢首评吧")
	}
	return nil
}

func (q *QueryCommentListFlow) packData() error {
	q.commentList = &CList{Comments: q.comments}
	return nil
}

func FillCommentListFields(comments *[]*models.Comment) error {
	if comments == nil {
		return errors.New("FillCommentListFields：传入comments指针为空")
	}
	commentsLen := len(*comments)
	if commentsLen == 0 {
		return errors.New("FillCommentListFields：传入comments无内容")
	}
	userInfoDAO := models.NewUserInfoDAO()
	for _, c := range *comments {
		_ = userInfoDAO.QueryUserInfoById(c.UserInfoId, &c.User)
		c.CreateDate = c.CreatedAt.Format("1-2")
	}
	return nil
}
