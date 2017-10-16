package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig("hosts.conf")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(conf)
}
