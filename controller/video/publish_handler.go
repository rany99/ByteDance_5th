package video

import (
	"ByteDance_5th/cover"
	"ByteDance_5th/models"
	"ByteDance_5th/server/video"
	"ByteDance_5th/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

var (
	videoSuffixMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	CoverSuffixMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

func PublishHandler(ctx *gin.Context) {
	rawId, _ := ctx.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		PublishError(ctx, "userid出错")
		return
	}
	title := ctx.PostForm("title")
	form, err := ctx.MultipartForm()
	if err != nil {
		PublishError(ctx, err.Error())
		return
	}

	files := form.File["data"]
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)
		log.Println("视频格式解析成功")
		if _, ok := videoSuffixMap[suffix]; !ok {
			PublishError(ctx, "不支持视频格式")
			continue
		}
		name := util.NewUnicFileName(userId)
		log.Println("视频名字生成成功")
		filename := name + suffix
		savePath := filepath.Join("./public", filename)
		log.Println("存储路径：", savePath)
		if err = ctx.SaveUploadedFile(file, savePath); err != nil {
			log.Println("视频保存失败")
			PublishError(ctx, err.Error())
			continue
		}
		snapShotFileName := util.NewUnicFileName(userId) + ".png"
		log.Println(snapShotFileName)
		snapShotPath := filepath.Join("./public", snapShotFileName)
		log.Println(snapShotPath)
		if err = cover.SnapShotFromVideo(savePath, snapShotPath, 1); err != nil {
			PublishError(ctx, err.Error())
			continue
		}
		if err = video.PostVideo(userId, filename, snapShotFileName, title); err != nil {
			log.Println("上传失败", err)
			PublishError(ctx, err.Error())
			continue
		}
		PublishOK(ctx, file.Filename+"上传成功")
	}
}

// PublishError 生成错误返回
func PublishError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// PublishOK 生成成功返回
func PublishOK(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}
