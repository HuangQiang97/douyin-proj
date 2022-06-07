package util

import (
	"douyin-proj/src/config"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := config.Init("../../../resources/application.ini"); err != nil {
		panic(m)
	}
	m.Run()
}

func TestReleaseToken(t *testing.T) {
	fmt.Println(jwtKey)
	token, err := ReleaseToken(uint(5))
	if err != nil {
		t.Errorf("releaseToen error = %v", err)
	} else {
		t.Log(token)
	}
}

func TestVerifyToken(t *testing.T) {
	originid := uint(3)
	token, err := ReleaseToken(originid)
	if err != nil {
		t.Errorf("releaseToen error = %v", err)
		return
	}
	parseid, err := VerifyToken(token)
	if err != nil {
		t.Errorf("verify faild = %v", err)
	} else if originid != parseid {
		t.Error("token is error")
	} else {
		t.Log(parseid)
	}
}

func TestVerifyTokenExpired(t *testing.T) {
	originid := uint(3)
	token, err := ReleaseToken(originid)
	if err != nil {
		t.Errorf("releaseToen error = %v", err)
	}
	time.Sleep(2 * time.Second)
	parseid, err := VerifyToken(token)
	if err != nil {
		t.Errorf("verify faild = %v", err)
	} else if originid != parseid {
		t.Error("token is error")
	} else {
		t.Log(parseid)
	}
}

func TestVerifyTokenIncorrect(t *testing.T) {
	originid := uint(3)
	token, err := ReleaseToken(originid)
	if err != nil {
		t.Errorf("releaseToen error = %v", err)
	}
	//修改token
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(token))
	temp := []byte(token)
	temp[index] = 'a'
	token = string(temp)

	parseid, err := VerifyToken(token)
	if err != nil {
		t.Errorf("verify faild = %v", err)
	} else if originid != parseid {
		t.Error("token is error")
	} else {
		t.Log(parseid)
	}
}
