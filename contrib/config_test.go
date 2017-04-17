package lamp

import (
	"testing"

	"github.com/yulefox/lamp/contrib"
)

func TestLoadConfig(t *testing.T) {
	conf := &contrib.LampConfig{}
	filename := "lamp.json"
	err := LoadConfig(filename, conf)

	if err != nil {
		t.Fatal(err)
	}
}
