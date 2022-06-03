package service

import (
	"douyin-proj/src/repository"
)

func CreateComment(userId uint, videoId uint, content string) (uint, error) {
	var comment = &repository.Comment{
		UserID:  userId,
		VideoID: videoId,
		Content: content,
	}

	err := repository.CreateComment(comment)
	if err != nil {
		return 0, err
	}

	return comment.ID, nil

}

func DeleteCommentById(id uint) error {
	// comment is not existed
	comment, err := repository.GetCommentById(id)
	if err != nil {
		return err
	}

	// comment delete failed
	if err := repository.DeleteComment(comment); err != nil {
		return err
	}

	//update video(commentCount)

	return nil

}
