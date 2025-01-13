package main

import (
	"fmt"
	"runtime"

	"github.com/chain5j/chain5j-pkg/cli"
	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"tg-keyword-reply-bot/pkg/config"
)

var (
	conf *config.Config[Config]

	bot   *api.BotAPI
	gcron *cron.Cron
)

func main() {
	cpuNum := runtime.NumCPU() // 获得当前设备的cpu核心数
	fmt.Println("CPU核心数:", cpuNum)
	runtime.GOMAXPROCS(cpuNum) // 设置需要用到的cpu数量
	initCli()
}

// 初始化命令行
func initCli() {
	rootCli := cli.NewCli(&cli.AppInfo{
		App:     "tgbot",
		Version: "tgbot",
		Welcome: "tgbot",
	})
	err := rootCli.InitFlags(true, func(rootFlags *pflag.FlagSet) {

	}, func(viper *viper.Viper) {
		// 这一步时,rootFlags已经获取到了值
		// 当config不为空时,才启用
		logViper := viper.Sub("log")
		logConfig := new(logger.LogConfig)
		err := logViper.Unmarshal(&logConfig)
		if err != nil {
			panic(err)
		}
		zap.InitWithConfig(logConfig)
	})
	if err != nil {
		logger.Fatal("initCli", "err", err)
	}

	rootCli.AddCommands(StartCommand(rootCli))
	rootCli.Execute()
}
