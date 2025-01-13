package main

type Config struct {
	Token       string `json:"token" mapstructure:"token" yaml:"token"`
	SuperUserId int64  `json:"super_user_id" mapstructure:"super_user_id" yaml:"super_user_id"`
	Debug       bool   `json:"debug" mapstructure:"debug" yaml:"debug"`
	ChanSize    int    `json:"chan_size" mapstructure:"chan_size" yaml:"chan_size"`
}
