package middleware

import (
	"bytes"
	"douyin-proj/src/config"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

// FfmpegTask 截图任务消息
type FfmpegTask struct {
	VideoName string
	ImageName string
}

// FfmpegTaskchan 截图任务消息信道
var FfmpegTaskchan chan FfmpegTask

// InitFfmpeg 初始化截图任务处理器
func InitFfmpeg() {
	FfmpegTaskchan = make(chan FfmpegTask, config.MaxMsgCount)
	for taskMsg := range FfmpegTaskchan {
		go func(f FfmpegTask) {
			err := GenCover(f.VideoName, f.ImageName, 1)
			if err != nil {
				FfmpegTaskchan <- f
				log.Fatal("派遣失败：重新派遣")
			}
			log.Printf("视频%v截图处理成功", f.VideoName)
		}(taskMsg)
	}
}

// GenCover 生成视频封面
func GenCover(videoName, imageName string, frameNum int) (err error) {
	videoPath := config.VideoSavePrefix + videoName
	imagePath := config.CoverSavePrefix + imageName
	// 抽取视频帧数据
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err
	}
	// 转图片数据
	img, err := imaging.Decode(buf)
	if err != nil {
		return err
	}
	// 保存图片
	if err = imaging.Save(img, imagePath); err != nil {
		return err
	}
	return nil
}
