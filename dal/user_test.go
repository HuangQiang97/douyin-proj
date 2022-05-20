package dal

import (
	"fmt"
	"testing"
)

func TestGetUserByName(t *testing.T){
	Init()

	users, err := GetUserByName("wmy")
	if err != nil{
		panic(err)
	}
	fmt.Println(len(users))
}