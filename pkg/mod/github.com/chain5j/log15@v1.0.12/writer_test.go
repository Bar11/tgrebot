// Package log
// 
// @author: xwc1125
// @date: 2021/3/29
package log

import (
	"log"
	"testing"
)

func TestTimeWriter_Write(t *testing.T) {
	timeWriter := &TimeWriter{
		Dir:        "./log",
		Compress:   true,
		ReserveDay: 30,
	}
	Info2 := log.New(timeWriter, " [Debug] ", log.LstdFlags)
	Info2.Println("this is debug")
}
