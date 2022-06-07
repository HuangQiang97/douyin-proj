package service

import (
	"bytes"
	"douyin-proj/src/config"
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"io/ioutil"
	"log"
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

// GetOutIp 获得本机外网IP
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

	//for _, a := range addrs {
	//	if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
	//		return ipnet.IP.String(), nil
	//	}
	//}
	for i := len(addrs) - 1; i >= 0; i-- {
		if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
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

func isVideo(suffix string) bool {
	videoTypes := []string{"avi", "wmv", "mpeg", "mp4", "m4v", "mov", "asf", "flv", "f4v", "rmvb", "rm", "3gp", "vob", "asx", "dat", "mkv", "webm", "3g2", "mpg", "mpe", "ts", "vob", "dat", "mkv", "lavf", "cpk", "dirac", "ram", "qt", "fli", "flc", "mod", "wmv", "avi", "dat", "asf", "mpeg", "mpg", " rm", "rmvb", "ram", "flv", "mp4", "3gp", " mov", "divx", "dv", "vob", "mkv", "qt", " cpk", "fli", "flc", "f4v", "m4v"}
	set := make(map[string]struct{}, len(videoTypes))
	for _, v := range videoTypes {
		set[v] = struct{}{}
	}
	_, ok := set[suffix]
	return ok
}

// SaveVideo 保存视频
func SaveVideo(file *multipart.FileHeader, uId uint, title string) (err error) {

	// 文件合法性判断
	if !isVideo(path.Ext(file.Filename)[1:]) {
		log.Printf("文件非视频文件。name:%s\n", file.Filename)
		return errors.New("文件非视频文件")
	}
	// 名称转16进制防止注入
	fileName := fmt.Sprintf("%x", file.Filename)
	// 截断防止过长数据库无法保存
	if len(fileName) > 36 {
		fileName = fileName[len(fileName)-36:]
	}
	fileId := fmt.Sprintf("%d_%d_%s", uId, time.Now().UnixNano(), fileName)
	videoPath := "./upload/video/" + fileId
	//保存视频
	if err = SaveFile(file, videoPath); err != nil {
		log.Printf("保存文件失败。path:%s,err:%s\n", videoPath, err)
		return err
	}

	// 获得本机ip
	ip, err := GetLocalIp()
	if err != nil {
		log.Printf("获取本机IP地址失败，err:%s\n", err)
		return err
	}
	// TODO
	ip = config.ServerConfig.HTTP_HOST

	videoUrl := "http://" + ip + ":" + config.ServerConfig.HTTP_PORT + videoPath[1:]
	//videoUrl := videoPath[1:]

	// 获得视频封面
	coverPath := "./upload/cover/" + strings.Split(fileId, ".")[0] + ".jpeg"
	if err = GenCover(videoPath, coverPath, 1); err != nil {
		log.Printf("获取封面失败。path:%s,err:%s\n", coverPath, err)
		return err
	}
	coverUrl := "http://" + ip + ":" + config.ServerConfig.HTTP_PORT + coverPath[1:]
	//coverUrl := coverPath[1:]

	// 剪切标题，防止过长数据库保存失败
	if len(title) > 255 {
		title = title[len(title)-255:]
	}

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
		log.Printf("保存视频信息失败。err:%s\n", err)
		return err
	}
	log.Printf("上传视频成功。videoPath:%s,coverPath:%s\n", videoPath, coverPath)

	return nil
}

// GetVideoList 获得用户视频列表
// authorId:要查询的用户,uid:当前登录用户
func GetVideoList(authorId, uid uint) (videoList []types.Video, err error) {
	// 获取用户信息
	var user *repository.User
	if util.UserExist(authorId) {
		user, err = util.GetUser(authorId)
	} else {
		user, err = repository.GetUserById(authorId)
		if user != nil {
			util.AddUser(user)
		}
	}
	//user, err := repository.GetUserById(authorId)

	if err != nil {
		log.Printf("根据uid获取用户信息失败。uid:%d,err:%s\n", authorId, err)
		return nil, err
	}
	// 获取用户发布视频列表
	videoIds, err := repository.GetVideoIdsByAuthorId(authorId)
	videos := make([]repository.Video, 0, len(videoIds))
	for _, vid := range videoIds {
		if util.VideoExist(uint(vid)) {
			video, _ := util.GetVideo(uint(vid))
			videos = append(videos, *video)
		} else {
			video, _ := repository.GetVideoById(uint(vid))
			videos = append(videos, *video)
			if video != nil {
				util.AddVideo(video)
			}
		}
	}

	//videos, err := repository.GetVideoByAuthorId(authorId)
	if err != nil {
		log.Printf("根据uid获取用发布视频列表失败。uid:%d,err:%s\n", authorId, err)
		return nil, err
	}

	// 填充用户信息
	var author = types.User{
		Id:            user.ID,
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow: repository.GetRelation(&repository.Relation{
			UserID:   uid,
			FollowID: authorId,
		}),
	}

	videoList = make([]types.Video, 0, len(videos))
	for _, v := range videos {
		// 查询当前用户是否给视频点赞过
		isFavorite := repository.IsFavorite(&repository.Favorite{
			UserID:  uid,
			VideoID: v.ID,
		})
		// 填充视频信息
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
	log.Printf("根据uid获取用户视频列表成功。uid:%d\n", authorId)
	return videoList, nil
}
