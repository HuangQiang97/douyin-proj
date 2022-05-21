package respository

import "testing"

func TestCreateRelation(t *testing.T) {
	r1 := Relation{
		UserID:   2,
		FollowID: 3,
	}
	if err := CreateRelation(&r1); err != nil {
		t.Error(err)
		return
	}
}

func TestDuplicationCreateRelation(t *testing.T) {
	r1 := Relation{
		UserID:   2,
		FollowID: 4,
	}
	if err := CreateRelation(&r1); err != nil {
		t.Error(err)
		return
	}
	t.Log("first insert success")
	r2 := Relation{
		UserID:   2,
		FollowID: 4,
	}
	if err := CreateRelation(&r2); err != nil {
		t.Error(err)
		return
	}
}
