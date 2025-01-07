package main

import (
	"github.com/chain5j/chain5j-pkg/cli"
	EventBus "github.com/chain5j/chain5j-pkg/eventbus"
	"github.com/chain5j/logger"
	"github.com/panjf2000/ants/v2"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"tg-keyword-reply-bot/db"
	"tg-keyword-reply-bot/pkg/config"
)

// startCommand 启动服务
type startCommand struct {
	log     logger.Logger
	rootCli *cli.Cli
	cmd     *cobra.Command

	msgBus EventBus.Bus
	stop   chan struct{}
}

// StartCommand 初始化cmd
func StartCommand(rootCli *cli.Cli) *cobra.Command {
	c := &startCommand{
		rootCli: rootCli,
		msgBus:  EventBus.New(),
		stop:    make(chan struct{}),
	}
	c.cmd = &cobra.Command{
		Use:   "start",
		Short: "Start tgbot",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.Log("tgbot")
			c.log = log
			configFile := rootCli.Viper().GetString("config")
			c.log.Info("config file path", "path", configFile)
			var err error
			conf, err = config.LoadConfigResource[Config](configFile, c.msgBus)
			if err != nil {
				return err
			}

			chPool1, err := ants.NewPool(conf.Config().ChanSize)
			if err != nil {
				log.Panic("new pool err", "err", err)
			}
			chPool = chPool1
			token := db.Init(conf.Config().Token)
			gcron = cron.New()
			gcron.Start()
			// 开始工作
			start(log, token)
			return nil
		},
	}
	c.addFlags()
	return c.cmd
}

// addFlags 添加flags
func (c *startCommand) addFlags() {
	// flags := c.rootCli.RootCmd().PersistentFlags()
	// // pprof
	// viper.BindPFlags(flags)
}
