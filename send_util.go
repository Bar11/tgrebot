package main

import (
	"github.com/chain5j/logger"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-keyword-reply-bot/db"
	"time"
)

// 发送文字消息
func sendMessage(log logger.Logger, msg api.MessageConfig) api.Message {
	if msg.Text == "" {
		log.Debug("message is nil")
		return api.Message{}
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send msg err", "err", err)
		return mmsg
	}
	db.AddMessageRecord(mmsg)
	// go deleteMessage(log, msg.ChatID, mmsg.MessageID)
	return mmsg
}

// 发送图片消息, 需要是已经存在的图片链接
func sendPhoto(log logger.Logger, chatId int64, filePath string) api.Message {
	msg := api.NewPhoto(chatId, api.FileURL(filePath))
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send photo err", "err", err)
	}
	db.AddMessageRecord(mmsg)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// sendGif 发送动图, 需要是已经存在的链接
func sendGif(log logger.Logger, chatId int64, filePath string) api.Message {
	file := api.NewDocument(chatId, api.FileURL(filePath))
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Error("bot send gif err", "err", err)
		return mmsg
	}
	db.AddMessageRecord(mmsg)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// sendVideo 发送视频, 需要是已经存在的视频连接
func sendVideo(log logger.Logger, chatId int64, filePath string) api.Message {
	file := api.NewVideo(chatId, api.FileURL(filePath))
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Error("bot send video err", "err", err)
		return mmsg
	}
	db.AddMessageRecord(mmsg)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// sendFile 发送文件, 必须是已经存在的文件链接
func sendFile(log logger.Logger, chatId int64, filePath string) api.Message {
	file := api.NewDocument(chatId, api.FilePath(filePath))
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Error("bot send file err", "err", err)
		return mmsg
	}
	db.AddMessageRecord(mmsg)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// deleteMessage 删除消息
func deleteMessage(log logger.Logger, gid int64, mid int) {
	time.Sleep(time.Second)
	_, _ = bot.Send(api.NewDeleteMessage(gid, mid))
}
