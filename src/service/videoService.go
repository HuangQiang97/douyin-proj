package service

import (
	"bytes"
	"douyin-proj/src/config"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
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

// GetLocalIp 获得本机IP
func GetOutIp() (ip string, err error) {
	// 获得地址
	response, err := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if err != nil {
		return "", err
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	ip = fmt.Sprintf("%s", string(body))
	return ip, nil
}

// LocalIPs return all non-loopback IPv4 addresses
func GetLocalIp() (string, error) {
	var ips string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", errors.New("can't get ip")
}

// GenCover 生成视频封面
func GenCover(videoPath, snapshotPath string, frameNum int) (err error) {
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
	if err = imaging.Save(img, snapshotPath); err != nil {
		return err
	}

	return nil
}

// SaveVideo 保存视频
func SaveVideo(file *multipart.FileHeader, uId uint, title string) (err error) {
	// 路径

	fileId := fmt.Sprintf("%d_%d_%s", uId, time.Now().Unix(), path.Clean(strings.Trim(file.Filename, "/")))
	videoPath := "./upload/video/" + fileId
	//保存视频
	if err = SaveFile(file, videoPath); err != nil {
		return err
	}
	// 获得本机ip

	ip, err := GetLocalIp()
	if err != nil {
		return err
	}
	// TODO
	ip = "10.192.18.194"

	videoUrl := "http://" + ip + ":" + config.ServerConfig.HTTP_PORT + videoPath[1:]
	//videoUrl := videoPath[1:]

	// 获得视频封面
	coverPath := "./upload/cover/" + strings.Split(fileId, ".")[0] + ".jpeg"
	if err = GenCover(videoPath, coverPath, 1); err != nil {
		return err
	}
	coverUrl := "http://" + ip + ":" + config.ServerConfig.HTTP_PORT + coverPath[1:]
	//coverUrl := coverPath[1:]

	// 保存视频数据到数据库
	video := repository.Video{
		AuthorID:      uId,
		PlayUrl:       videoUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		CreatedAt:     uint64(time.Now().Unix()),
	}
	if err = repository.CreateVideo(&video); err != nil {
		return err
	}
	return nil
}

// GetVideoList 获得用户视频列表
func GetVideoList(authorId uint) (videoList []types.Video, err error) {
	user, err := repository.GetUserById(authorId)
	if err != nil {
		return nil, err
	}
	videos, err := repository.GetVideoByAuthorId(authorId)
	if err != nil {
		return nil, err
	}

	var author = types.User{
		Id:            user.ID,
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow:      false,
	}
	videoList = make([]types.Video, 0, len(videos))
	for _, v := range videos {
		isFavorite := repository.IsFavorite(&repository.Favorite{
			UserID:  authorId,
			VideoID: authorId,
		})
		videoList = append(videoList, types.Video{
			Id:            v.ID,
			Author:        author,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		})
	}
	return videoList, nil
}
