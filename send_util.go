package main

import (
	"github.com/chain5j/logger"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"tg-keyword-reply-bot/db"
	"time"
)

// 发送文字消息
func sendMessage(log logger.Logger, msg api.MessageConfig) api.Message {
	//fmt.Println(msg.Text)
	if msg.Text == "" {
		log.Debug("message is nil")
		return api.Message{}
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send msg err", "err", err)
		return mmsg
	}

	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	// go deleteMessage(log, msg.ChatID, mmsg.MessageID)
	return mmsg
}

func SendKeyboardButtonData(log logger.Logger, gid int64, msg api.Message) api.Message {
	filebuttonCallBackData := "option1"
	photobuttonCallBackData := "option2"
	gifbuttonCallBackData := "option3"
	videobuttonCallBackData := "option4"
	voicebuttonCallBackData := "option5"
	filebutton := api.InlineKeyboardButton{
		Text:         "upload file",
		CallbackData: &filebuttonCallBackData,
	}
	photobutton := api.InlineKeyboardButton{
		Text:         "upload photo",
		CallbackData: &photobuttonCallBackData,
	}
	gifbutton := api.InlineKeyboardButton{
		Text:         "upload gif",
		CallbackData: &gifbuttonCallBackData,
	}
	videobutton := api.InlineKeyboardButton{
		Text:         "upload video",
		CallbackData: &videobuttonCallBackData,
	}
	voicebutton := api.InlineKeyboardButton{
		Text:         "upload voice",
		CallbackData: &voicebuttonCallBackData,
	}
	keyboard := api.InlineKeyboardMarkup{
		InlineKeyboard: [][]api.InlineKeyboardButton{
			{filebutton},
			api.NewInlineKeyboardRow(gifbutton, photobutton),
			api.NewInlineKeyboardRow(voicebutton, videobutton),
		}}
	mmsg := api.NewMessage(msg.Chat.ID, msg.Text)
	mmsg.ReplyMarkup = keyboard
	mmmsg, err := bot.Send(mmsg)
	if err != nil {
		log.Error("bot send msg err", "err", err)
	}

	//keyboard2 := api.ReplyKeyboardRemove{
	//	RemoveKeyboard: true,
	//}
	//mmsg.ReplyMarkup = keyboard2
	//mmmsg, _ := bot.Send(mmsg)

	return mmmsg

}

// 发送图片消息, 需要是已经存在的图片链接
func sendPhoto(log logger.Logger, chatId int64, filePath string) api.Message {
	var msg api.Chattable
	if strings.HasPrefix(filePath, "http") {
		msg = api.NewPhoto(chatId, api.FileURL(filePath))
	} else {
		msg = api.NewPhoto(chatId, api.FileID(filePath))
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send photo err", "err", err)
	}
	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// 发送图片消息, 需要是已经存在的图片链接
func sendVoice(log logger.Logger, chatId int64, filePath string) api.Message {
	var msg api.Chattable
	if strings.HasPrefix(filePath, "http") {
		msg = api.NewVoice(chatId, api.FileURL(filePath))
	} else {
		msg = api.NewVoice(chatId, api.FileID(filePath))
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send voice err", "err", err)
		return mmsg
	}
	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	return mmsg
}

// sendGif 发送动图, 需要是已经存在的链接
func sendGif(log logger.Logger, chatId int64, filePath string) api.Message {
	var msg api.Chattable
	if strings.HasPrefix(filePath, "http") {
		msg = api.NewDocument(chatId, api.FileURL(filePath))
	} else {
		msg = api.NewDocument(chatId, api.FileID(filePath))
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send gif err", "err", err)
		return mmsg
	}
	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// sendVideo 发送视频, 需要是已经存在的视频连接
func sendVideo(log logger.Logger, chatId int64, filePath string) api.Message {
	var msg api.Chattable
	if strings.HasPrefix(filePath, "http") {
		msg = api.NewVideo(chatId, api.FileURL(filePath))
	} else {
		msg = api.NewVideo(chatId, api.FileID(filePath))
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send video err", "err", err)
		return mmsg
	}
	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// sendFile 发送文件, 必须是已经存在的文件链接
func sendFile(log logger.Logger, chatId int64, filePath string) api.Message {
	var msg api.Chattable
	if strings.HasPrefix(filePath, "http") {
		msg = api.NewDocument(chatId, api.FileURL(filePath))
	} else {
		msg = api.NewDocument(chatId, api.FileID(filePath))
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Error("bot send file err", "err", err)
		return mmsg
	}
	url, messageType := GetUrlFromServer(mmsg, bot)
	db.AddMessageRecord(mmsg, url, messageType)
	// go deleteMessage(log, chatId, mmsg.MessageID)
	return mmsg
}

// deleteMessage 删除消息
func deleteMessage(log logger.Logger, gid int64, mid int) {
	time.Sleep(time.Second)
	_, _ = bot.Send(api.NewDeleteMessage(gid, mid))
}
