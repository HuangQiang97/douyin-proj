package config

import "testing"

func TestInit(t *testing.T) {
	if err := Init("../../resources/application.ini"); err != nil {
		t.Errorf("Init error = %v", err)
	}
	t.Log(MySQLConfig)
	t.Log(ServerConfig)
	t.Log(SecretKey)
	t.Logf("test finish!")
}
