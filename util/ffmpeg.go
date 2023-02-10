package util

import "C"
import (
	"ByteDance_5th/config"
	"errors"
	"fmt"
	"log"
	"unsafe"
)

type VideoToCover struct {
	InputPath  string
	OutputPath string
	StartTime  string
	KeepTime   string
	Filter     string
	FrameCount int64
	debug      bool
}

var videoToCover VideoToCover

func NewVideoToCover() *VideoToCover {
	return &videoToCover
}

var videoChanger VideoToCover

const (
	inputVideoPathOption = "-i"
	startTimeOption      = "-ss"
	keepTimeOption       = "-t"
	videoFilterOption    = "-vf"
	formatToImageOption  = "-f"
	autoReWriteOption    = "-y"
	framesOption         = "-frames:v"
)

var (
	DefaultVideoSuffix = ".mp4"
	DefaultImageSuffix = ".jpg"
)

func (v *VideoToCover) Debug() {
	v.debug = true
}

func (v *VideoToCover) GetQueryString() (ret string, err error) {
	if v.InputPath == "" {
		err = errors.New("输入路径为空")
	}
	if v.OutputPath == "" {
		err = errors.New("输出路径为空")
	}
	ret = config.Conf.Path.FfmpegPath
	ret += ParamAndParam(inputVideoPathOption, v.InputPath)
	ret += ParamAndParam(formatToImageOption, "cover_of_")
	if v.Filter != "" {
		ret += ParamAndParam(videoFilterOption, v.Filter)
	}
	if v.StartTime != "" {
		ret += ParamAndParam(startTimeOption, v.StartTime)
	}
	if v.KeepTime != "" {
		ret += ParamAndParam(keepTimeOption, v.KeepTime)
	}
	if v.FrameCount != 0 {
		ret += ParamAndParam(framesOption, fmt.Sprintf("%d", v.FrameCount))
	}
	ret += ParamAndParam(autoReWriteOption, v.OutputPath)
	return
}

// ParamAndParam 拼接ffmpeg参数信息
func ParamAndParam(p1, p2 string) string {
	return fmt.Sprintf(" %s %s ", p1, p2)
}

// GetDefaultVideoSuffix 返回默认图片类型
func GetDefaultVideoSuffix() string {
	return DefaultVideoSuffix
}

// GetDefaultImageSuffix 返回默认图片类型
func GetDefaultImageSuffix() string {
	return DefaultImageSuffix
}

func (v *VideoToCover) ExecCommand(cmd string) error {
	if v.debug {
		log.Println(cmd)
	}
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	status := C.startCmd(cCmd)
	if status != 0 {
		return errors.New("视频切截图失败")
	}
	return nil
}
