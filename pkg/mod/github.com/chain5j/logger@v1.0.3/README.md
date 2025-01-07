# logger

本项目通过logger接口的适配，让log15、zap、logrus的命令终端输出及文件保存格式一致。

## 特性
- 统一终端输出效果
- 终端输出及文件保存配置分开
- 终端可颜色高亮输出
- 可根据module选择性输出

## 项目地址
Github: https://github.com/chain5j/logger

## 使用

```go
import (
	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap"
	//"github.com/chain5j/logger/logrus"
	//"github.com/chain5j/logger/log15"
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	log := zap.InitWithConfig(&logger.LogConfig{
		Console: logger.ConsoleLogConfig{
			Level:    4,
			Modules:  "*",
			ShowPath: false,
			UseColor: true,
			Console:  true,
		},
		File: logger.FileLogConfig{
			Level:    4,
			Modules:  "*",
			Save:     true,
			FilePath: "./logs",
			FileName: "log.json",
		},
	})

	var wg sync.WaitGroup
	startTime := logger.CurrentTime()

	debugLog := log.New("Debug")
	infoLog := log.New("Info")
	errorLog := log.New("Error")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			infoLog.Info("=========================", "i", i)
			for j := 0; j < 100; j++ {
				debugLog.Debug("test1 debug", "i", i, "j", j)
				if i%9 == 0 {
					infoLog.Info("test2 info", "i", i, "j", j)
				}
				if i%13 == 0 {
					errorLog.Error("test2 info", "i", i, "j", j)
				}
			}
		}(i)
	}
	wg.Wait()
	log.Info("总耗时", "elapsed", logger.CurrentTime()-startTime)
}
```

## LICENSE
`logger` 的源码允许用户在遵循 [Apache 2.0 开源证书](LICENSE) 规则的前提下使用。

## 版权
Copyright@2022 chain5j

![chain5j](./chain5j.png)