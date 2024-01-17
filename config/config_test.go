package config

import "testing"

func TestInit(t *testing.T) {
	Init("../config/config.yaml")
	t.Log(conf)
}
