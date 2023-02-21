package cache

import (
	"ByteDance_5th/pkg/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

// redis 初始化
func init() {
	rdb = redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Conf.RD.IP, config.Conf.RD.Port),
			Password: "",
			DB:       config.Conf.RD.Database,
		})
}

// ProxyCache 缓存层
type ProxyCache struct {
}

// 代理层
var (
	proxyCacheOperation ProxyCache
)

func NewProxyIndexMap() *ProxyCache {
	return &proxyCacheOperation
}

// GetVideoFavor 获取点赞状态 ret： true点赞 false未点赞
func (p *ProxyCache) GetVideoFavor(uid, vid int64) bool {
	key := fmt.Sprintf("30322favorite:%d", uid)
	return rdb.SIsMember(ctx, key, vid).Val()
}

// SetVideoFavor isFavor: true点赞 false取消点赞
func (p *ProxyCache) SetVideoFavor(uid, vid int64, isFavor bool) {
	key := fmt.Sprintf("30322favorite:%d", uid)
	if isFavor {
		rdb.SAdd(ctx, key, vid)
		return
	}
	rdb.SRem(ctx, key, vid)
}

// GetAFollowB 判断A是否关注了B
func (p *ProxyCache) GetAFollowB(a, b int64) bool {
	key := fmt.Sprintf("30322followrelation:%d", a)
	return rdb.SIsMember(ctx, key, b).Val()
}

// SetAFollowB isFollowed：true已关注 false未关注
func (p *ProxyCache) SetAFollowB(a, b int64, isFollowed bool) {
	key := fmt.Sprintf("30322followrelation:%d", a)
	if isFollowed {
		rdb.SAdd(ctx, key, b)
		return
	}
	rdb.SRem(ctx, key, b)
}
