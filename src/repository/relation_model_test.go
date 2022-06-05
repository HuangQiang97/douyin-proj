package repository

import (
	"fmt"
	"testing"
)

func TestCreateRelation(t *testing.T) {
	if err := CreateRelationWithCount(7, 3); err != nil {
		t.Error(err)
		return
	}
}

//func TestDuplicationCreateRelation(t *testing.T) {
//	r1 := Relation{
//		UserID:   2,
//		FollowID: 4,
//	}
//	if err := CreateRelation(&r1); err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log("first insert success")
//	r2 := Relation{
//		UserID:   2,
//		FollowID: 4,
//	}
//	if err := CreateRelation(&r2); err != nil {
//		t.Error(err)
//		return
//	}
//}

func TestDeleteRelationWithCount(t *testing.T) {
	if err := DeleteRelationWithCount(7, 3); err != nil {
		t.Error(err)
		return
	}
}

func TestGetFollow(t *testing.T) {
	users, err := GetFollow(20)
	if err != nil {
		t.Error(err)
		return
	}
	for _, u := range users {
		fmt.Println(u)
	}
}

func TestGetFans(t *testing.T) {
	users, err := GetFans(20)
	if err != nil {
		t.Error(err)
		return
	}
	for _, u := range users {
		fmt.Println(u)
	}
}
