package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/video"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FavoriteListResponse struct {
	models.CommonResponse
	*video.FavoriteList
}

type ProxyFavoriteListController struct {
	uid int64
	*gin.Context
}

func QueryFavoriteListController(ctx *gin.Context) {
	NewProxyFavoriteListController(ctx).Do()
}

func NewProxyFavoriteListController(ctx *gin.Context) *ProxyFavoriteListController {
	return &ProxyFavoriteListController{
		Context: ctx,
	}
}

func (p *ProxyFavoriteListController) Do() {
	if err := p.ParseJson(); err != nil {
		p.SendFailed(err.Error())
	}
	favoriteList, err := video.QueryFavoriteList(p.uid)
	if err != nil {
		p.SendFailed(err.Error())
	}
	p.SendSucceed(favoriteList)
}

func (p *ProxyFavoriteListController) ParseJson() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New("uid解析错误")
	}
	p.uid = uid
	return nil
}

func (p *ProxyFavoriteListController) SendFailed(msg string) {
	p.JSON(http.StatusOK, FavoriteListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FavoriteList: nil,
	})
}

func (p *ProxyFavoriteListController) SendSucceed(favoriteList *video.FavoriteList) {
	p.JSON(http.StatusOK, FavoriteListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		FavoriteList: favoriteList,
	})
}
