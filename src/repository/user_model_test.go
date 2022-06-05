package repository

import (
	"douyin-proj/src/config"
	"douyin-proj/src/database"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	if err := config.Init("../../resources/application.ini"); err != nil {
		panic(err)
	}
	if err := database.Init(); err != nil {
		panic(err)
	}
	DB = database.MySQLDb
	fmt.Printf("db= %v\n", DB)
	m.Run()
}

func TestCreateUser(t *testing.T) {
	var user = &User{UserName: "qaz5", Password: "qwer4"}
	err := CreateUser(user)
	if err != nil {
		t.Errorf("create user failed : %v", err)
	}
}

func TestGetUserById(t *testing.T) {
	var id uint = 2
	user, err := GetUserById(id)
	if err != nil {
		t.Errorf("get user by id=%v error =%v", id, err)
		return
	}
	fmt.Println(user)
}

func TestGetUsersByIds(t *testing.T) {
	var ids = []uint{1, 2, 3}
	users, err := GetUsersByIds(ids)
	if err != nil {
		t.Errorf("get user by ids=%v error =%v", ids, err)
		return
	}
	for _, user := range users {
		fmt.Println(user)
	}
}

func TestGetUserByName(t *testing.T) {
	user, err := GetUserByName("qaz6")
	if err != nil {
		t.Errorf("get user by name error = %v", err)
		return
	}
	fmt.Println(user)
}

//func TestUpdateFollowAndFans(t *testing.T) {
//	if err := UpdateFollowAndFans(3, 4, -1); err != nil {
//		t.Errorf("update user follow&fans error = %v", err)
//		return
//	}
//}
//
//func TestUpdateFollow(t *testing.T) {
//	if err := UpdateFollow(3, -1); err != nil {
//		t.Errorf("update user follow&fans error = %v", err)
//		return
//	}
//}
//
//func TestUpdateFollow2(t *testing.T) {
//	M := 100000
//	w := sync.WaitGroup{}
//	w.Add(M)
//	for i := 0; i < 100000; i++ {
//		go func() {
//			UpdateFollow(14, 1)
//			w.Done()
//		}()
//	}
//	w.Wait()
//
//}

func TestGetUserResponse(t *testing.T) {
	var uid, id = uint(3), uint(2)
	user, isFolllow := GetUserResponse(uid, id)
	fmt.Println(user, isFolllow)

}
