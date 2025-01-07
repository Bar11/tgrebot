// Package logger
//
// @author: xwc1125
package logger

import "testing"

func TestLogger(t *testing.T) {
	Printf("%s", "111")
	Println("222", "111")
	Info("222", "111", "111")
	Fatal("111")
	Println("333", "444")
}
