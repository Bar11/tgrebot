package config

import (
	"sync"

	"github.com/chain5j/chain5j-pkg/cli"
	EventBus "github.com/chain5j/chain5j-pkg/eventbus"
	"github.com/chain5j/logger"
	"github.com/fsnotify/fsnotify"
	"tg-keyword-reply-bot/pkg/compare"
)

// Config 模块
// 加载本地静态配置文件
// 校验必填参数
// 组装参数
// 提供参数查询接口
type Config[T any] struct {
	log         logger.Logger
	localConfig *sync.Pool // 节点本地启动配置,*models.LocalConfig
	msgBus      EventBus.Bus
}

// LoadConfigResource 加载本地配置文件
func LoadConfigResource[T any](configFile string, msgBus EventBus.Bus) (*Config[T], error) {
	// 1、读取静态配置 yaml
	// 2、校验静态文件URL是否可用，必填配置参数是否存在，若不存在抛出异常
	// 3、将静态配置文件内容组装到config结构体中： ChainConfig LocalConfig
	c := &Config[T]{
		log:    logger.Log("config"),
		msgBus: msgBus,
	}
	var localConfig T
	c.localConfig = &sync.Pool{
		New: func() interface{} {
			return localConfig
		},
	}

	err := cli.LoadConfigWithEvent(configFile, "tgbot", &localConfig, func(e fsnotify.Event) {
		c.log.Info("Config module watch file changed", "eventName", e.Name)
		oldConfig := c.Config()
		var localConfig1 T
		if err := cli.LoadConfig(configFile, "", &localConfig1); err == nil {
			c.localConfig.Put(localConfig1)
			c.log.Info("Local config update success")
			c.compareConfig(oldConfig, localConfig1)
		} else {
			c.log.Error("load local config err", "err", err)
		}
	})
	if err != nil {
		c.log.Error("load config err", "err", err)
		return nil, err
	}

	//zap.InitWithConfig(&localConfig.Log)
	return c, err
}

func (c *Config[T]) compareConfig(oldConfig, newConfig T) {
	_, diffCount := compare.CompareStruct(oldConfig, newConfig)
	if len(diffCount) == 0 {
		return
	}

	for key, _ := range diffCount {
		c.msgBus.Publish(key)
	}
}

func (c *Config[T]) WatchConfig(subKey string, fn func()) error {
	return c.msgBus.Subscribe(subKey, fn)
}

// Config ...
func (c *Config[T]) Config() T {
	return c.localConfig.Get().(T)
}
