// Package config
package config

import (
	"fmt"
	"sync"
	"testing"

	EventBus "github.com/chain5j/chain5j-pkg/eventbus"
	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
)

type AA struct {
	Log logger.LogConfig `json:"log" mapstructure:"log"`
	A   string           `json:"a" mapstructure:"a" yaml:"a"`
	B   string           `json:"b" mapstructure:"b" yaml:"b"`
}

func TestLoadConfigResource(t *testing.T) {
	msgBus := EventBus.New()
	conf, err := LoadConfigResource[AA]("config_test.yml", msgBus)
	if err != nil {
		t.Fatal(err)
	}
	aa := conf.Config()
	zap.InitWithConfig(&aa.Log)
	conf.WatchConfig("A", func() {
		fmt.Println("配置已更新，自行获取最新配置")
	})
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
