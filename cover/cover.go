//----------!!!Attention:请确保ffmpeg.exe已经置于GoPath路径下!!!----------

package cover

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

//videoPath 视频保存路径
//snapShotPath 截图保存路径
//frameNum 截图帧数

// SnapShotFromVideo 生成截图
func SnapShotFromVideo(videoPath, snapShotPath string, frameNum int) (err error) {

	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}

	err = imaging.Save(img, snapShotPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err
	}
	log.Println("截图成功")
	return nil
}
