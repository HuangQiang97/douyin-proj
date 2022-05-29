package repository

import "testing"

func TestCreateFavorite(t *testing.T) {
	f := Favorite{
		UserID:  3,
		VideoID: 7,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
}

func TestDuplicateCreateFavorite(t *testing.T) {
	f := Favorite{
		UserID:  3,
		VideoID: 8,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("first insert favorite success")
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
}

func TestUndoFavorite(t *testing.T) {
	f := Favorite{
		UserID:  3,
		VideoID: 9,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("insert favorite success")
	if err := UndoFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("undo favorite success")
}