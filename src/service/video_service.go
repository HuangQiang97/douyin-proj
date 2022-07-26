package service

import (
	"douyin-proj/src/config"
	"douyin-proj/src/repository"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/types"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// SaveFile 保存文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

// isVideo 判断文件是否是视频
func isVideo(suffix string) bool {
	videoTypes := []string{"avi", "wmv", "mpeg", "mp4", "m4v", "mov", "asf", "flv", "f4v", "rmvb", "rm", "3gp", "vob", "asx", "dat", "mkv", "webm", "3g2", "mpg", "mpe", "ts", "vob", "dat", "mkv", "lavf", "cpk", "dirac", "ram", "qt", "fli", "flc", "mod", "wmv", "avi", "dat", "asf", "mpeg", "mpg", " rm", "rmvb", "ram", "flv", "mp4", "3gp", " mov", "divx", "dv", "vob", "mkv", "qt", " cpk", "fli", "flc", "f4v", "m4v"}
	set := make(map[string]struct{},len(videoTypes))
	for _, v := range videoTypes {
		set[v] = struct{}{}
	}
	_, ok := set[suffix]
	return ok
}

// SaveVideo 保存视频
func SaveVideo(file *multipart.FileHeader, userId uint, title string) (err error) {
	// 文件合法性判断
	if !isVideo(path.Ext(file.Filename)[1:]) {
		log.Printf("文件非视频文件。name:%s\n", file.Filename)
		return errors.New("文件非视频文件")
	}
	// 生成视频文件名防止注入
	videoName := uuid.NewV4().String() + path.Ext(file.Filename)
	videoPath := config.VideoSavePrefix + videoName
	//保存视频
	if err = SaveFile(file, videoPath); err != nil {
		log.Printf("保存文件失败。path:%s,err:%s\n", videoPath, err)
		return err
	}
	//生成视频文件名防止注入
	imageName := uuid.NewV4().String() + ".png"
	//向队列中添加生成截屏任务
	middleware.FfmpegTaskchan <- middleware.FfmpegTask{
		VideoName: videoName,
		ImageName: imageName,
	}
	// 保存视频数据到数据库
	video := repository.Video{
		AuthorID:      userId,
		PlayUrl:       config.PlayUrlPrefix + videoName,
		CoverUrl:      config.CoverUrlPrefix + imageName,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		CreatedAt:     uint64(time.Now().Unix()),
	}
	if err = repository.CreateVideo(&video); err != nil {
		log.Printf("保存视频信息失败。err:%s\n", err)
		return err
	}
	// 添加新增视频到缓存
	if middleware.ExistUserPublish(userId) {
		middleware.AddUserPublish(userId, video.ID)
	}
	return nil
}

// GetVideoList 获得用户视频列表
func GetVideoList(authorId, currId uint) (videoList []types.Video, err error) {
	// 获取用户发布视频列表
	var videoIds []uint
	// 尝试从缓存获取
	if middleware.ExistUserPublish(authorId) {
		videoIds, _ = middleware.GetUserPublish(authorId)
	} else {
		videoIds, _ = repository.GetVideoIdsByAuthorId(authorId)
		for _, videoId := range videoIds {
			middleware.AddUserPublish(authorId, videoId)
		}
	}
	// go程
	var wg sync.WaitGroup
	wg.Add(len(videoIds))
	// 填充视频信息
	videoList = make([]types.Video, len(videoIds), len(videoIds))
	for idx, videoId := range videoIds {
		go func(idx int, videoId uint) {
			basicVideo, _ := repository.GetVideoById(videoId)
			user, _ := GetUserInfo(basicVideo.AuthorID, currId)
			video := types.Video{
				Id:            basicVideo.ID,
				PlayUrl:       basicVideo.PlayUrl,
				CoverUrl:      basicVideo.CoverUrl,
				FavoriteCount: uint64(getVideoFavoriteCount(videoId)),
				CommentCount:  uint64(getVideoCommentCount(videoId)),
				Title:         basicVideo.Title,
				Author:        *user,
				IsFavorite:    isFavorite(videoId, currId),
			}
			videoList[idx] = video
			wg.Done()
		}(idx, videoId)
	}
	wg.Wait()
	log.Printf("根据uid获取用户视频列表成功。currId:%d\n", authorId)
	return videoList, nil
}
