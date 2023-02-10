package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListResponse struct {
	models.CommonResponse
	*video.PublishList
}

// ProxyQueryVideoList 防火层
type ProxyQueryVideoList struct {
	ctx *gin.Context
}

func NewProxyQueryVideoList(ctx *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{ctx: ctx}
}

// QueryVideoListController Controller层
func QueryVideoListController(ctx *gin.Context) {
	p := NewProxyQueryVideoList(ctx)
	rawId, _ := ctx.Get("user_id")
	uid, ok := rawId.(int64)
	if !ok {
		p.QueryVideoListFailed("uid解析错误")
	}
	if err := p.DoQueryVideoListByUid(uid); err != nil {
		p.QueryVideoListFailed(err.Error())
	}
}

// DoQueryVideoListByUid 使用uid知行查询
func (p *ProxyQueryVideoList) DoQueryVideoListByUid(uid int64) error {
	videoList, err := video.QueryPublishListByUid(uid)
	if err != nil {
		return err
	}
	p.QueryVideoListOk(videoList)
	return nil
}

// QueryVideoListOk 获取成功
func (p *ProxyQueryVideoList) QueryVideoListOk(list *video.PublishList) {
	p.ctx.JSON(http.StatusOK, ListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		PublishList: list,
	})
}

// QueryVideoListFailed 获取失败
func (p *ProxyQueryVideoList) QueryVideoListFailed(msg string) {
	p.ctx.JSON(http.StatusOK, ListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		PublishList: nil,
	})
}
