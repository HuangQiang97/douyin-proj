package dal

import (
	"fmt"
	"testing"
)

func TestRelation(t *testing.T){
	Init()
	r := &Relation{
		UserID: 1,
		ToUserID: 3,
	}
	InsertRelation(r)
	r2 := &Relation{
		UserID: 1,
		ToUserID: 4,
	}
	InsertRelation(r2)
	user,_ := GetAllFollower(3)
	if len(user) != 0{
		t.Error(" failed")
	}
	fmt.Println(user)
}

