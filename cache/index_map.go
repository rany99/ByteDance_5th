package cache

import (
	"ByteDance_5th/pkg/common"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ProxyIndexMap struct {
}

var (
	ctx = context.Background()
	rdb *redis.Client
)

var proxyIndexOperation ProxyIndexMap

func init() {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", common.Conf.RD.IP, common.Conf.RD.Port),
			Password: "",
			DB:       common.Conf.RD.Database,
		})
}

func NewProxyIndexMap() *ProxyIndexMap {
	return &proxyIndexOperation
}

// GetVideoFavor 获取点赞状态 ret： true点赞 false未点赞
func (p *ProxyIndexMap) GetVideoFavor(userid, videoid int64) bool {
	key := fmt.Sprintf("favor:%d", userid)
	return rdb.SIsMember(ctx, key, videoid).Val()
}

// SetVideoFavor isFavor: true点赞 false取消点赞
func (p *ProxyIndexMap) SetVideoFavor(userid, videoId int64, isFavor bool) {
	key := fmt.Sprintf("favor:%d", userid)
	if isFavor {
		rdb.SAdd(ctx, key, videoId)
		return
	}
	rdb.SRem(ctx, key, videoId)
}

// GetAFollowB 判断A是否关注了B
func (p *ProxyIndexMap) GetAFollowB(a, b int64) bool {
	key := fmt.Sprintf("relation:%d", a)
	return rdb.SIsMember(ctx, key, b).Val()
}

// SetAFollowB isFollowed：true已关注 false未关注
func (p *ProxyIndexMap) SetAFollowB(a, b int64, isFollowed bool) {
	key := fmt.Sprintf("relation:%d", a)
	if isFollowed {
		rdb.SAdd(ctx, key, b)
	}
	rdb.SRem(ctx, key, b)
}
