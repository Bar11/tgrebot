// Package logger
//
// @author: xwc1125
package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"reflect"
	"strings"
	"time"
)

func GetWriter(fileConfig FileLogConfig) (io.WriteCloser, error) {
	// WithMaxAge和WithRotationCount 二者只能设置一个，
	// WithMaxAge 设置文件清理前的最长保存时间，
	// WithRotationCount 设置文件清理前最多保存的个数。
	// filename是指向最新日志的链接
	hook, err := rotatelogs.New(
		fileConfig.GetLogFile()+".%Y%m%d%H",
		rotatelogs.WithLinkName(fileConfig.GetLogFile()),                                   // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*time.Duration(fileConfig.GetMaxAge())),             // 文件最大保存时间,保存30天
		rotatelogs.WithRotationTime(time.Hour*time.Duration(fileConfig.GetRotationTime())), // 日志切割时间间隔,切割频率 24小时
	)
	if err != nil {
		return nil, err
	}
	return hook, nil
}

func ColorRender(str string, color int, weight int, extraArgs ...interface{}) string {
	//闪烁效果
	isBlink := int64(0)
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}
	//下划线效果
	isUnderLine := int64(0)
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}

func GetColorByLevel(level Lvl) int {
	switch level {
	case LvlFatal:
		return 35
	case LvlError:
		return 31
	case LvlWarn:
		return 33
	case LvlInfo:
		return 32
	case LvlDebug:
		return 36
	case LvlTrace:
		return 34
	default:
		return 0
	}
}

// CurrentTime 返回毫秒
func CurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}