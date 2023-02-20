package constantval

const (
	CreateCommentActionType int64 = 1
	DeleteCommentActionType int64 = 2
)

const (
	SendMsgActionType int64 = 1
)

const (
	FollowSucceed             int = 0
	FollowUserNoExist         int = 1
	PostFollowActionTypeWrong int = 2
	CantFollowSelf            int = 3
	FollowAgain               int = 4
)
