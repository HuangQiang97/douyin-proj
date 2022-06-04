package repository

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVideosByIds(t *testing.T) {
	var ids = []uint{1, 8}
	videos, err := GetVideosByIds(ids)
	if err != nil {
		t.Errorf("get Video by ids error =%v", err)
		return
	}
	for _, v := range videos {
		fmt.Println(v)
	}
}

func TestGetVideoByAuthorId(t *testing.T) {
	var id = uint(7)
	videos, err := GetVideoByAuthorId(id)
	if err != nil {
		t.Errorf("get Video by author_id error =%v", err)
		return
	}
	for _, v := range videos {
		fmt.Println(v)
	}
}

func TestCreateVideo(t *testing.T) {
	video := Video{
		AuthorID:      3,
		PlayUrl:       "qsaasdd",
		CoverUrl:      "qazxx",
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         "gsz video 02",
		CreatedAt:     uint64(time.Now().Unix()),
	}
	if err := CreateVideo(&video); err != nil {
		t.Errorf("create video failed ,error= %v", err)
		return
	}
}

func TestGetVideoByAuthorIdWithFavorite(t *testing.T) {
	var author_id, id = uint(3), uint(5)
	videolist := GetVideoByAuthorIdWithFavorite(author_id, id)
	for _, v := range videolist {
		fmt.Println(v)
	}
}
