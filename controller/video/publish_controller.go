package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/video"
	"ByteDance_5th/util"
	"ByteDance_5th/util/information"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func PublishController(ctx *gin.Context) {

	// 解析userId
	rawId, _ := ctx.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		PublishFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 解析题目
	title := ctx.PostForm("title")
	form, err := ctx.MultipartForm()
	if err != nil {
		PublishFailed(ctx, errortype.ParseTitleErr)
		return
	}

	// 提取文件
	files := form.File["data"]
	for _, file := range files {

		//视频解析
		suffix := filepath.Ext(file.Filename)
		if _, ok := admitVideoSuffixMap[suffix]; !ok {
			msg := errortype.WrongVideoTypeErr + suffix
			PublishFailed(ctx, msg)
			continue
		}
		//log.Println("视频格式解析成功")

		// 使用雪花算法生成唯一视频id（非数据库id）
		node, _ := util.NewWorker(1)
		randomId := node.GetId()
		fileID := fmt.Sprintf("%v", randomId)

		//生成视频存储路径
		filename := fileID + suffix
		savePath := filepath.Join("./public/video", filename)
		if err = ctx.SaveUploadedFile(file, savePath); err != nil {
			PublishFailed(ctx, errortype.VideoSaveErr)
			continue
		}

		//生成封面存储路径
		snapShotFileName := fileID + ".png"
		snapShotPath := filepath.Join("./public/cover", snapShotFileName)

		//生成视频封面
		if err = information.SnapShotFromVideo(savePath, snapShotPath, 1); err != nil {
			PublishFailed(ctx, err.Error())
			continue
		}

		//上传视频
		if err = video.PostVideo(userId, filename, snapShotFileName, title); err != nil {
			//log.Println("上传失败", err)
			PublishFailed(ctx, err.Error())
			continue
		}

		models.NewUserInfoDAO().WorkCntAddOneByUid(userId)
		
		PublishSucceed(ctx, file.Filename+"上传成功")
	}
}

// PublishFailed 生成错误返回
func PublishFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// PublishSucceed 生成成功返回
func PublishSucceed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  msg,
	})
}

var (
	admitVideoSuffixMap = map[string]bool{
		".mp4":  true,
		".avi":  true,
		".wmv":  true,
		".flv":  true,
		".mpeg": true,
		".mov":  true,
	}
)
