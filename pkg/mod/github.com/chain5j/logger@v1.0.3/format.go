// Package logger
//
// @author: xwc1125
package logger

import "time"

const (
	TermMsgJust = 44

	DefaultTimestampFormat     = time.RFC3339
	DefaultTermTimestampFormat = "2006-01-02 15:04:05.000"
	FieldKeyMsg                = "msg"
	FieldKeyLevel              = "lvl"
	FieldKeyTime               = "t"
	FieldKeyFile               = "f"
	FieldKeyModule             = "module"
	FieldKeyError              = "err"
)
