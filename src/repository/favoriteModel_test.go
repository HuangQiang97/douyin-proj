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

func TestListFavorite(t *testing.T) {
	f := Favorite{
		UserID:  3,
		VideoID: 10,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("insert favorite #1 success")

	f = Favorite{
		UserID:  3,
		VideoID: 11,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("insert favorite #2 success")

	videoIds, err := GetFavoriteVideoIdsByUserId(3)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("List favorite success", videoIds)
}

func TestListFavoriteEmpty(t *testing.T) {
	videoIds, err := GetFavoriteVideoIdsByUserId(1843)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("List empty favorite success", videoIds)
}

func TestIsFavorite(t *testing.T) {
	f := Favorite{
		UserID:  3,
		VideoID: 12,
	}
	if err := CreateFavorite(&f); err != nil {
		t.Error(err)
		return
	}
	t.Log("insert favorite success")

	result := IsFavorite(&f)
	t.Log("Is favorite: ", result)

	f = Favorite{
		UserID:  2032,
		VideoID: 0,
	}
	result = IsFavorite(&f)
	t.Log("Is favorite: ", result)
}
