package information

import (
	"ByteDance_5th/pkg/common"
	"fmt"
	"math/rand"
	"strconv"
)

const BackgroundCnt int = 6

func GetBackGroundUrl() string {
	i := rand.Intn(100)
	fileName := strconv.Itoa(i%BackgroundCnt) + ".jpg"
	var url string = fmt.Sprintf("http://%s:%d/static/background/%s", common.Conf.SE.IP, common.Conf.SE.Port, fileName)
	return url
}
