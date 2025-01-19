package main

import (
	"strconv"
	"time"
	"unicode"

	"github.com/chain5j/logger"
	"tg-keyword-reply-bot/common"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkSuperUser(log logger.Logger, user api.User) bool {
	if conf.Config().SuperUserId == user.ID {
		log.Info("user is super user", "uid", user.ID)
		return true
	} else {
		log.Info("user is not super user", "uid", user.ID)
		return false
	}
}

// 检查是否是群组的管理员
func checkAdmin(log logger.Logger, gid int64, user api.User) bool {

	admins, _ := bot.GetChatAdministrators(api.ChatAdministratorsConfig{ChatConfig: api.ChatConfig{ChatID: gid, SuperGroupUsername: ""}})
	uid := user.ID
	if conf.Config().SuperUserId > 0 && uid == conf.Config().SuperUserId {
		log.Info("user is super user", "uid", uid)
		return true
	}
	for _, user := range admins {
		if uid == user.User.ID {
			log.Info("user is group admin", "uid", uid)
			return true
		}
	}
	punish(log, gid, user)
	return false
}

// checkInGroup 检查是不是新加的群或者新开的人
func checkInGroup(id int64) bool {
	for _, gid := range common.AllGroupId {
		if gid == id {
			return true
		}
	}
	return false
}

// punish 惩罚
func punish(log logger.Logger, gid int64, user api.User) {
	botme, _ := bot.GetChatMember(api.GetChatMemberConfig{ChatConfigWithUser: api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID}})
	msg := api.NewMessage(gid, "")
	if botme.CanRestrictMembers {
		// 禁言60秒
		banMember(log, gid, user.ID, 60)
		msg.Text = "[" + user.String() + "](tg://user?id=" + strconv.FormatInt(user.ID, 10) + ")乱玩管理员命令,禁言一分钟"
		msg.ParseMode = "Markdown"
	} else {
		msg.Text = "[" + user.String() + "](tg://user?id=" + strconv.FormatInt(user.ID, 10) + ")不要乱玩管理员命令"
		msg.ParseMode = "Markdown"
	}
	sendMessage(log, msg)
}

// 禁言群员
func banMember(log logger.Logger, gid int64, uid int64, sec int64) {
	if sec <= 0 {
		sec = 9999999999999
	}
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	chatPermissions := &api.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
	}
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        time.Now().Unix() + sec,
		Permissions:      chatPermissions,
	}
	_, _ = bot.Send(restricconfig)
}

// 解除禁言
func unbanMember(log logger.Logger, gid int64, uid int64) {
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	chatPermissions := &api.ChatPermissions{
		CanSendMessages:       true,
		CanSendMediaMessages:  true,
		CanSendOtherMessages:  true,
		CanAddWebPagePreviews: true,
	}
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        9999999999999,
		Permissions:      chatPermissions,
	}
	_, _ = bot.Send(restricconfig)
}

// 踢出群员
func kickMember(log logger.Logger, gid int64, uid int64) {
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	chatPermissions := &api.ChatPermissions{
		CanSendMessages:       false,
		CanSendMediaMessages:  false,
		CanSendPolls:          false,
		CanSendOtherMessages:  false,
		CanAddWebPagePreviews: false,
		CanChangeInfo:         false,
		CanInviteUsers:        false,
		CanPinMessages:        false,
	}
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        9999999999999,
		Permissions:      chatPermissions,
	}
	_, _ = bot.Send(restricconfig)
}

// 解除禁止
func unkickMember(log logger.Logger, gid int64, uid int64) {
	chatuserconfig := api.ChatMemberConfig{ChatID: gid, UserID: uid}
	chatPermissions := &api.ChatPermissions{
		CanSendMessages:       true,
		CanSendMediaMessages:  true,
		CanSendPolls:          true,
		CanSendOtherMessages:  true,
		CanAddWebPagePreviews: true,
		CanChangeInfo:         true,
		CanInviteUsers:        true,
		CanPinMessages:        true,
	}
	restricconfig := api.RestrictChatMemberConfig{
		ChatMemberConfig: chatuserconfig,
		UntilDate:        time.Now().Unix(),
		Permissions:      chatPermissions,
	}
	_, _ = bot.Send(restricconfig)
}

// 返回群组的所有管理员, 用来进行一次性@
func getAdmins(log logger.Logger, gid int64) string {
	admins, _ := bot.GetChatAdministrators(api.ChatAdministratorsConfig{ChatConfig: api.ChatConfig{ChatID: gid}})
	list := ""
	for _, admin := range admins {
		user := admin.User
		if user.IsBot {
			continue
		}
		list += "[" + user.String() + "](tg://user?id=" + strconv.FormatInt(admin.User.ID, 10) + ")\r\n"
	}
	return list
}

// checkQingzhen 检查文字中是否包含阿拉伯文
func checkQingzhen(text string) bool {
	for _, c := range text {
		if unicode.Is(unicode.Scripts["Arabic"], c) {
			return true
		}
	}
	return false
}
