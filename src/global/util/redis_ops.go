package util

import (
	"douyin-proj/src/database"
	"douyin-proj/src/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func AddVideo(video *repository.Video) error {
	if video == nil {
		return errors.New("视频为空")
	}
	base := fmt.Sprintf("video_%d", video.ID)
	metaData, err := json.Marshal(video)
	if err != nil {
		log.Printf("序列化视频异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "meta", metaData).Err(); err != nil {
		log.Printf("保存视频异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "favo", video.FavoriteCount).Err(); err != nil {
		log.Printf("保存视频点赞数异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "comm", video.CommentCount).Err(); err != nil {
		log.Printf("保存视频评论数异常。err:%s\n", err)
		return err
	}
	return nil
}

func VideoFavoIncr(id uint) error {
	base := fmt.Sprintf("video_%d", id)
	if err := database.RedisClient.HIncrBy(base, "favo", 1).Err(); err != nil {
		log.Printf("视频点赞数加一异常。err:%s\n", err)
		return err
	}
	return nil
}

func VideoFavoDecr(id uint) error {
	base := fmt.Sprintf("video_%d", id)
	if err := database.RedisClient.HIncrBy(base, "favo", -1).Err(); err != nil {
		log.Printf("视频点赞数减一异常。err:%s\n", err)
		return err
	}
	return nil
}

func VideoCommIncr(id uint) error {
	base := fmt.Sprintf("video_%d", id)
	if err := database.RedisClient.HIncrBy(base, "comm", 1).Err(); err != nil {
		log.Printf("视频评论数加一异常。err:%s\n", err)
		return err
	}
	return nil
}

func VideoCommDecr(id uint) error {
	base := fmt.Sprintf("video_%d", id)
	if err := database.RedisClient.HIncrBy(base, "comm", -1).Err(); err != nil {
		log.Printf("视频评论数数减一异常。err:%s\n", err)
		return err
	}
	return nil
}

func AddUser(user *repository.User) error {
	if user == nil {
		return errors.New("用户为空")
	}
	base := fmt.Sprintf("user_%d", user.ID)
	metaData, err := json.Marshal(user)
	if err != nil {
		log.Printf("序列化用户异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "meta", metaData).Err(); err != nil {
		log.Printf("保存用户信息异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "follow", user.FollowCount).Err(); err != nil {
		log.Printf("保存用户关注信息异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, "follower", user.FansCount).Err(); err != nil {
		log.Printf("保存用户粉丝信息异常。err:%s\n", err)
		return err
	}
	return nil
}

func UserFollowIncr(id uint) error {
	base := fmt.Sprintf("user_%d", id)
	if err := database.RedisClient.HIncrBy(base, "follow", 1).Err(); err != nil {
		log.Printf("用户关注数加一异常。err:%s\n", err)
		return err
	}
	return nil
}

func UserFollowDecr(id uint) error {
	base := fmt.Sprintf("user_%d", id)
	if err := database.RedisClient.HIncrBy(base, "follow", -1).Err(); err != nil {
		log.Printf("用户关注数减一异常。err:%s\n", err)
		return err
	}
	return nil
}

func UserFollowerIncr(id uint) error {
	base := fmt.Sprintf("user_%d", id)
	if err := database.RedisClient.HIncrBy(base, "follower", 1).Err(); err != nil {
		log.Printf("用户粉丝数加一异常。err:%s\n", err)
		return err
	}
	return nil
}

func UserFollowerDecr(id uint) error {
	base := fmt.Sprintf("user_%d", id)
	if err := database.RedisClient.HIncrBy(base, "follower", -1).Err(); err != nil {
		log.Printf("用户粉丝数一异常。err:%s\n", err)
		return err
	}
	return nil
}

func GetVideo(id uint) (video *repository.Video, err error) {
	base := fmt.Sprintf("video_%d", id)
	metaData, err := database.RedisClient.HGet(base, "meta").Bytes()
	if err != nil {
		log.Printf("获取视频异常。err:%s\n", err)
		return nil, err
	}
	err = json.Unmarshal(metaData, &video)
	if err != nil {
		log.Printf("反序列化视频异常。err:%s\n", err)
		return nil, err
	}
	video.FavoriteCount, err = database.RedisClient.HGet(base, "favo").Uint64()
	if err != nil {
		log.Printf("获取视频点赞个数异常。err:%s\n", err)
		return nil, err
	}
	video.CommentCount, err = database.RedisClient.HGet(base, "comm").Uint64()
	if err != nil {
		log.Printf("获取视频评论个数异常。err:%s\n", err)
		return nil, err
	}
	return video, nil
}

func GetUser(id uint) (user *repository.User, err error) {
	base := fmt.Sprintf("user_%d", id)
	metaData, err := database.RedisClient.HGet(base, "meta").Bytes()
	if err != nil {
		log.Printf("获取用户信息异常。err:%s\n", err)
		return nil, err
	}
	err = json.Unmarshal(metaData, &user)
	if err != nil {
		log.Printf("反序列化用户异常。err:%s\n", err)
		return nil, err
	}
	user.FollowCount, err = database.RedisClient.HGet(base, "follow").Uint64()
	if err != nil {
		log.Printf("获取用户关注个数异常。err:%s\n", err)
		return nil, err
	}
	user.FansCount, err = database.RedisClient.HGet(base, "follower").Uint64()
	if err != nil {
		log.Printf("获取用户粉丝个数异常。err:%s\n", err)
		return nil, err
	}
	return user, nil
}

func AddComments(videoId uint, comments *[]repository.Comment) (err error) {
	if comments == nil {
		return errors.New("评论集合为空")
	}
	base := fmt.Sprintf("comm_%d", videoId)
	for _, comm := range *comments {
		data, err := json.Marshal(comm)

		if err != nil {
			log.Printf("序列化评论异常。err:%s\n", err)
			return err
		}
		if err = database.RedisClient.HSet(base, strconv.Itoa(int(comm.ID)), data).Err(); err != nil {
			log.Printf("保存评论异常。err:%s\n", err)
			return err
		}
	}
	return nil
}

//AddComment 保存评论，视频评论数加一
func AddComment(videoId uint, comment *repository.Comment) (err error) {
	if comment == nil {
		return errors.New("评论为空")
	}
	base := fmt.Sprintf("comm_%d", videoId)
	data, err := json.Marshal(comment)
	if err != nil {
		log.Printf("序列化评论异常。err:%s\n", err)
		return err
	}
	if err = database.RedisClient.HSet(base, strconv.Itoa(int(comment.ID)), data).Err(); err != nil {
		log.Printf("保存评论异常。err:%s\n", err)
		return err
	}
	return nil
}

//DeleteComment  删除评论，视频评论数-1
func DeleteComment(videoId uint, commId uint) (err error) {
	base := fmt.Sprintf("comm_%d", videoId)
	if err = database.RedisClient.HDel(base, strconv.Itoa(int(commId))).Err(); err != nil {
		log.Printf("删除评论异常。err:%s\n", err)
		return err
	}
	return nil
}

//func DeleteComments(videoId uint) (err error) {
//	base := fmt.Sprintf("comm_%d", videoId)
//	if err = database.RedisClient.Del(base).Err(); err != nil {
//		log.Printf("删除评论异常。err:%s\n", err)
//		return err
//	}
//	return nil
//}

func GetComments(videoId uint) ([]repository.Comment, error) {
	base := fmt.Sprintf("comm_%d", videoId)
	re, err := database.RedisClient.HGetAll(base).Result()
	if err != nil {
		log.Printf("获取评论异常。err:%s\n", err)
		return nil, err
	}
	comments := make([]repository.Comment, 0, len(re))
	for _, data := range re {
		comm := repository.Comment{}
		err = json.Unmarshal([]byte(data), &comm)
		if err != nil {
			log.Printf("反序列化评论异常。err:%s\n", err)
			return nil, err
		}
		comments = append(comments, comm)
	}
	return comments, nil
}

func VideoExist(id uint) bool {
	base := fmt.Sprintf("video_%d", id)
	re, _ := database.RedisClient.Exists(base).Result()
	return re > 0
}
func UserExist(id uint) bool {
	base := fmt.Sprintf("user_%d", id)
	re, _ := database.RedisClient.Exists(base).Result()
	return re > 0
}

func CommentsExist(id uint) bool {
	base := fmt.Sprintf("comm_%d", id)
	re, _ := database.RedisClient.Exists(base).Result()
	return re > 0
}
