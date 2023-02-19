package information

import (
	"ByteDance_5th/pkg/common"
	"fmt"
	"math/rand"
	"strconv"
)

const AvatarCnt int = 16

// GetAvatarUrl 生成头像url
// 由于本次客户端中并没有给出相应的用于上传头像的接口，因此在public/avatar文件中预存了16张图片用作头像
func GetAvatarUrl() string {
	i := rand.Intn(100)
	fileName := strconv.Itoa(i%AvatarCnt) + ".jpg"
	var url string = fmt.Sprintf("http://%s:%d/static/avatar/%s", common.Conf.SE.IP, common.Conf.SE.Port, fileName)
	return url
}
