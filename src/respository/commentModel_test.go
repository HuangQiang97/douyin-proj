package respository

import (
	"testing"
)

func TestCreateComment(t *testing.T) {
	comment := &Comment{
		UserID:  3,
		VideoID: 6,
		Content: "TEST",
	}
	if err := CreateComment(comment); err != nil {
		t.Errorf("create comment error = %v", err)
		return
	}
	t.Log(comment)
}
