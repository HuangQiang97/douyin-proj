package middleware

import (
	"douyin-proj/src/config"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

// InitRedis 初始化redis连接
func InitRedis() {
	// 根据redis配置初始化一个客户端
	db, _ := strconv.Atoi(config.RedisConfig.DB)
	addr := fmt.Sprintf("%s:%s", config.RedisConfig.Addr, config.RedisConfig.Port)
	password := config.RedisConfig.Password
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,     // redis地址
		Password: password, // redis密码，没有则留空
		DB:       db,       // 默认数据库，默认是0
	})
}

///////////////////////////////////////////////////////////////////////////////////////////

// addToSet 向集合添加元素
func addToSet(key string, value uint) error {

	err := RedisClient.SAdd(key, value).Err()
	if err != nil {
		return err
	}
	return RedisClient.Expire(key, config.RedisExpireTime).Err()

}

// delFromSet 从集合中删除元素
func delFromSet(key string, value uint) error {
	RedisClient.Expire(key, config.RedisExpireTime)
	return RedisClient.SRem(key, value).Err()
}

// countFromSet 统计集合中元素个数
func countFromSet(key string) (uint, error) {
	RedisClient.Expire(key, config.RedisExpireTime)
	cnt, err := RedisClient.SCard(key).Result()
	return uint(cnt), err
}

// getIdsFromSet 获取集合中成员
func getIdsFromSet(key string) ([]uint, error) {
	RedisClient.Expire(key, config.RedisExpireTime)
	strIds, err := RedisClient.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]uint, len(strIds), len(strIds))

	for idx, strId := range strIds {
		id, _ := strconv.ParseInt(strId, 10, 64)
		ids[idx] = uint(id)
	}
	return ids, nil
}

///////////////////////////////////用户发表视频记录操作////////////////////////////////////////////////////////

func AddUserPublish(userId uint, videoId uint) error {
	key := fmt.Sprintf("publish_user_%d", userId)
	return addToSet(key, videoId)
}

func GetUserPublish(id uint) ([]uint, error) {
	key := fmt.Sprintf("publish_user_%d", id)
	return getIdsFromSet(key)
}

func ExistUserPublish(id uint) bool {
	key := fmt.Sprintf("publish_user_%d", id)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt > 0
}

///////////////////////////////////用户关注他人记录操作////////////////////////////////////////////////////////

func AddUserFollowing(idCurr uint, idFollow uint) error {
	key := fmt.Sprintf("following_user_%d", idCurr)
	return addToSet(key, idFollow)
}

func DelUserFollowing(idCurr uint, idFollow uint) error {
	key := fmt.Sprintf("following_user_%d", idCurr)
	return delFromSet(key, idFollow)
}

func CountUserFollowing(idCurr uint) (uint, error) {
	key := fmt.Sprintf("following_user_%d", idCurr)
	return countFromSet(key)
}

func GetUserFollowing(idCurr uint) ([]uint, error) {
	key := fmt.Sprintf("following_user_%d", idCurr)
	return getIdsFromSet(key)
}

func ExistUserFollowing(id uint) bool {
	key := fmt.Sprintf("following_user_%d", id)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt > 0
}

func ExistFollowRelation(currId uint, idFollow uint) bool {
	key := fmt.Sprintf("following_user_%d", currId)
	ans, _ := RedisClient.SIsMember(key, idFollow).Result()
	return ans
}

//////////////////////////////////用户被他人关注记录操作/////////////////////////////////////////////////////////

func AddUserFollower(idCurr uint, idFollower uint) error {
	key := fmt.Sprintf("follower_user_%d", idCurr)
	return addToSet(key, idFollower)
}

func DelUserFollower(idCurr uint, idFollower uint) error {
	key := fmt.Sprintf("follower_user_%d", idCurr)
	return delFromSet(key, idFollower)
}

func CountUserFollower(idCurr uint) (uint, error) {
	key := fmt.Sprintf("follower_user_%d", idCurr)
	return countFromSet(key)
}

func GetUserFollower(idCurr uint) ([]uint, error) {
	key := fmt.Sprintf("follower_user_%d", idCurr)
	return getIdsFromSet(key)
}

func ExistUserFollower(id uint) bool {
	key := fmt.Sprintf("follower_user_%d", id)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt > 0
}

//////////////////////////////////////用户点赞视频记录操作/////////////////////////////////////////////////////

func AddUserFavorite(userId uint, videoId uint) error {
	key := fmt.Sprintf("favorite_user_%d", userId)
	return addToSet(key, videoId)
}

func DelUserFavorite(userId uint, videoId uint) error {
	key := fmt.Sprintf("favorite_user_%d", userId)
	return delFromSet(key, videoId)
}

func GetUserFavorite(userId uint) ([]uint, error) {
	key := fmt.Sprintf("favorite_user_%d", userId)
	return getIdsFromSet(key)
}

func ExistUserFavorite(userId uint) bool {
	key := fmt.Sprintf("favorite_user_%d", userId)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt > 0
}

func ExistFavoriteRelation(userId uint, videoId uint) bool {
	key := fmt.Sprintf("favorite_user_%d", userId)
	ans, _ := RedisClient.SIsMember(key, videoId).Result()
	return ans
}

//////////////////////////////////////用户视频被点赞数统计操作/////////////////////////////////////////////////////
func InitVideoFavoriteCount(videoId uint, count int64) error {
	key := fmt.Sprintf("favorite_video_count_%d", videoId)
	return RedisClient.Set(key, count, config.RedisExpireTime).Err()
}

func IncrVideoFavoriteCount(videoId uint) (int64, error) {
	key := fmt.Sprintf("favorite_video_count_%d", videoId)
	return RedisClient.Incr(key).Result()
}

func DecrVideoFavoriteCount(videoId uint) (int64, error) {
	key := fmt.Sprintf("favorite_video_count_%d", videoId)
	return RedisClient.Decr(key).Result()
}

func GetVideoFavoriteCount(videoId uint) (uint, error) {
	key := fmt.Sprintf("favorite_video_count_%d", videoId)
	ans, err := RedisClient.Get(key).Result()
	count, _ := strconv.ParseInt(ans, 10, 64)
	return uint(count), err
}

func ExistVideoFavoriteCount(videoId uint) bool {
	key := fmt.Sprintf("favorite_video_count_%d", videoId)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt == 1
}

/////////////////////////////////////视频评论记录操作//////////////////////////////////////////////////////

func AddVideoComment(videoId uint, commentId uint) error {
	key := fmt.Sprintf("comment_video_%d", videoId)
	return addToSet(key, commentId)
}

func DelVideoComment(videoId uint, commentId uint) error {
	key := fmt.Sprintf("comment_video_%d", videoId)
	return delFromSet(key, commentId)
}

func CountVideoComment(videoId uint) (uint, error) {
	key := fmt.Sprintf("comment_video_%d", videoId)
	return countFromSet(key)
}

func ExistVideoComment(id uint) bool {
	key := fmt.Sprintf("comment_video_%d", id)
	cnt, _ := RedisClient.Exists(key).Result()
	return cnt > 0
}

func GetVideoComment(videoId uint) ([]uint, error) {
	key := fmt.Sprintf("comment_video_%d", videoId)
	return getIdsFromSet(key)
}

///////////////////////////////////////////////////////////////////////////////////////////
