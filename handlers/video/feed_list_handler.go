package video

//import (
//	"ByteDance_5th/models"
//	"ByteDance_5th/server/video"
//	"github.com/gin-gonic/gin"
//)
//
//type FeedResponse struct {
//	models.CommonResponse
//	*video.FeedList
//}
//
//func FeedListHandler(ctx *gin.Context) {
//	p := NewProxyFeedList(ctx)
//	token, ok := ctx.Query("token")
//	if !ok {
//		err := p.DoNoToken();
//	}
//
//}
//
//// ProxyFeedList f防火层
//type ProxyFeedList struct {
//	*gin.Context
//}
//
//func NewProxyFeedList(ctx *gin.Context) *ProxyFeedList {
//	return &ProxyFeedList{ctx}
//}
//
//func (p *ProxyFeedList) DoNoToken() error {
//	timeStamp := p.Query("lstset_time")
//	var  =
//}
