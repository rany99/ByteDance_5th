package video

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/server/video"
	"ByteDance_5th/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

var (
	videoSuffixMap = map[string]bool{
		".mp4":  true,
		".avi":  true,
		".wmv":  true,
		".flv":  true,
		".mpeg": true,
		".mov":  true,
	}
)

func PublishHandler(ctx *gin.Context) {
	rawId, _ := ctx.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		PublishError(ctx, errortype.ParseUserIdErr)
		return
	}
	title := ctx.PostForm("title")
	form, err := ctx.MultipartForm()
	if err != nil {
		PublishError(ctx, errortype.ParseTitleErr)
		return
	}

	files := form.File["data"]
	for _, file := range files {
		//视频解析
		suffix := filepath.Ext(file.Filename)
		if _, ok := videoSuffixMap[suffix]; !ok {
			msg := errortype.WrongVideoTypeErr + suffix
			PublishError(ctx, msg)
			continue
		}
		//log.Println("视频格式解析成功")

		//生成视频存储路径
		name := util.NewUnicFileName(userId)
		filename := name + suffix
		savePath := filepath.Join("./public/video", filename)
		//log.Println("存储路径：", savePath)
		if err = ctx.SaveUploadedFile(file, savePath); err != nil {
			PublishError(ctx, errortype.VideoSaveErr)
			continue
		}

		//生成封面存储路径
		snapShotFileName := util.NewUnicFileName(userId) + ".png"
		//log.Println(snapShotFileName)
		snapShotPath := filepath.Join("./public/cover", snapShotFileName)
		//log.Println(snapShotPath)

		//生成视频封面
		if err = util.SnapShotFromVideo(savePath, snapShotPath, 1); err != nil {
			PublishError(ctx, err.Error())
			continue
		}

		//上传视频
		if err = video.PostVideo(userId, filename, snapShotFileName, title); err != nil {
			//log.Println("上传失败", err)
			PublishError(ctx, err.Error())
			continue
		}
		PublishOK(ctx, file.Filename+"上传成功")
	}
}

// PublishError 生成错误返回
func PublishError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// PublishOK 生成成功返回
func PublishOK(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}
