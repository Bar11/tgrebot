package main

import (
	"github.com/panjf2000/ants/v2"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/chain5j/logger"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-keyword-reply-bot/common"
	"tg-keyword-reply-bot/db"
)

var (
	chPool *ants.Pool
)

// start 开始工作
func start(log logger.Logger, botToken string) {
	var err error

	bot, err = api.NewBotAPI(botToken)
	if err != nil {
		log.Error("new bot api err", "err", err)
		log.Fatal(err)
	}
	bot.Debug = conf.Config().Debug
	log.Info("authorized on account", "bot_name", bot.Self.UserName, "bot_id", bot.Self.ID)
	if conf.Config().SuperUserId != 0 {
		log.Info("super user is has been set", "user_id", conf.Config().SuperUserId)
	}
	u := api.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error("bot get updates err", "err", err)
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		chPool.Submit(func() {
			processUpdate(log, &update)
		})
	}
}

// processUpdate 对于每一个update的单独处理
func processUpdate(log logger.Logger, update *api.Update) {
	upmsg := update.Message
	//jsonString, _ := json.Marshal(upmsg)
	//fmt.Println(string(jsonString))
	url, messageType := GetUrlFromServer(*upmsg, bot)
	db.AddMessageRecord(*upmsg, url, messageType)
	log.Debug("update msg", "msg", upmsg.Text)
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	log.Debug("message from update", "gid", gid, "uid", uid, "mid", upmsg.MessageID, "uname", upmsg.From.String())
	// 检查是不是新开的群或者新加的人
	in := checkInGroup(gid)
	if !in {
		// 不在就需要加入, 内存中加一份, 数据库中添加一条空规则记录
		var title string
		TypeChat := upmsg.Chat.Type
		if TypeChat == "private" {
			title = upmsg.From.FirstName + " " + upmsg.From.LastName
		} else {
			title = upmsg.Chat.Title
		}
		db.AddNewGroup(gid, title, TypeChat)
		common.AddNewGroup(gid)
		log.Info("add new group", "gid", gid)
	}
	// 判断msg是否是命令
	if upmsg.IsCommand() {
		chPool.Submit(func() {
			// 处理指令
			processCommand(log, update)
		})
	} else {
		chPool.Submit(func() {
			processReplyCommand(log, update)
			processReply(log, update)
		})
		// 新用户通过用户名检查是否是清真
		if upmsg.NewChatMembers != nil {
			// todo 新加入的用户，可以推送欢迎语
			for _, auser := range upmsg.NewChatMembers {
				if checkQingzhen(auser.UserName) ||
					checkQingzhen(auser.FirstName) ||
					checkQingzhen(auser.LastName) {
					banMember(log, gid, uid, -1)
				}
			}
		}
		// 检查清真并剔除
		if checkQingzhen(upmsg.Text) {
			_ = api.NewDeleteMessage(gid, upmsg.MessageID)
			banMember(log, gid, uid, -1)
		}
	}
}

func processReply(log logger.Logger, update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	// 根据请求的内容查询回复的内容
	replyText := findKey(gid, upmsg.Text)
	if replyText == "delete" {
		_ = api.NewDeleteMessage(gid, upmsg.MessageID)
	} else if strings.HasPrefix(replyText, "ban:") {
		sec, err := strconv.ParseInt(replyText[4:], 10, 64)
		if err != nil {
			log.Error("parse int error", "err", err)
		}
		msg = api.NewMessage(gid, "")
		if -1 < sec && sec < 9999999999999 {
			msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") Speech involving sensitive words, banned for " + strconv.FormatInt(sec, 10) + " second.If you have any questions, please contact the administrator."
		} else {
			msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") Speech involving sensitive words, permanently banned.If you have any questions, please contact the administrator"

		}
		msg.ParseMode = "Markdown"
		sendMessage(log, msg)
		banMember(log, gid, uid, sec)
		deleteMessage(log, gid, upmsg.MessageID)
	} else if strings.HasPrefix(replyText, "photo:") {
		sendPhoto(log, gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "voice:") {
		sendVoice(log, gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "gif:") {
		sendGif(log, gid, replyText[4:])
	} else if strings.HasPrefix(replyText, "video:") {
		sendVideo(log, gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "document:") {
		sendFile(log, gid, replyText[9:])
	} else if replyText != "" {
		msg = api.NewMessage(gid, replyText)
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = upmsg.MessageID
		sendMessage(log, msg)
	}
}

// processCommand 指令处理
func processCommand(log logger.Logger, update *api.Update) {
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	msg := api.NewMessage(update.Message.Chat.ID, "")
	// 上传发送的本地文件（包含图片、语音、视频、文档）
	if upmsg.ReplyToMessage != nil {
		switch upmsg.Command() {
		case "local":
			if checkAdmin(log, gid, *upmsg.From) {
				order := upmsg.CommandArguments()
				if order != "" {
					url, messageType := GetUrlFromServer(*upmsg.ReplyToMessage, bot)
					if url != "" {
						switch messageType {
						case "Photo":
							Photo := upmsg.ReplyToMessage.Photo[len(upmsg.ReplyToMessage.Photo)-1]
							order += "===photo:" + Photo.FileID
						case "Video":
							order += "===video:" + upmsg.ReplyToMessage.Video.FileID
						case "Document":
							order += "===document:" + upmsg.ReplyToMessage.Document.FileID
						case "Voice":
							order += "===voice:" + upmsg.ReplyToMessage.Voice.FileID
						default:
							return
						}
					} else {
						return
					}
					if checkSuperUser(log, *upmsg.From) {
						for _, gid_ := range common.AllGroupId {
							addRule(gid_, order)
							msg.Text = "所有规则添加成功: " + order
						}
					} else {
						addRule(gid, order)
						msg.Text = "规则添加成功: " + order
					}

				} else {
					msg.Text = localText
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true
				}
				sendMessage(log, msg)
			}
		}

	}
	switch upmsg.Command() {
	case "start":
		msg.Text = "本机器人能够自动回复特定关键词"
		sendMessage(log, msg)
	case "help":
		msg.Text = helpText
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		sendMessage(log, msg)
	case "sensitive":
		order := upmsg.CommandArguments()
		if checkSuperUser(log, *upmsg.From) {
			if order != "" {
				for _, gid := range common.AllGroupId {
					addBanRule(gid, order)
					msg.Text = "敏感词规则添加成功: " + order
				}
			} else {
				msg.Text = addBanText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
		} else if checkAdmin(log, gid, *upmsg.From) {
			if order != "" {
				addBanRule(gid, order)
				msg.Text = "敏感词规则添加成功: " + order
			} else {
				msg.Text = addBanText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
		} else {
			return
		}
		sendMessage(log, msg)

	case "allChat":
		if checkSuperUser(log, *upmsg.From) && upmsg.Chat.Type == "private" {
			rules, err := db.GetAllGroup(log)
			if err != nil {
				log.Error("GetAllGroup DB Error")
			} else {
				for _, r := range rules {
					msg.Text += r.ChatTitle + ", " + strconv.FormatInt(r.GroupId, 10) + ", " + r.ChatType + "\r\n"
				}
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
		} else {
			msg.Text = "need SuperUser and private chat with bot."
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			//deleteMessage(log, gid, upmsg.MessageID)
		}
		sendMessage(log, msg)
	case "add":
		if checkAdmin(log, gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" && strings.Contains(order, "===") {
				addRule(gid, order)
				msg.Text = "规则添加成功: " + order
			} else {
				msg.Text = addText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessage(log, msg)
		}
	case "addForAll":
		if checkSuperUser(log, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" && strings.Contains(order, "===") {
				for _, gid := range common.AllGroupId {
					addRule(gid, order)
				}
				msg.Text = "为所有群组、好友，添加规则成功: " + order
			} else {
				msg.Text = addForAllText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessage(log, msg)
		} else {
			msg.Text = "该命令仅适用超级管理员"
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(log, msg)
		}
	case "copy":
		if checkSuperUser(log, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				order = strings.Replace(order, " ", "", -1)
				//新增 copy to 逻辑
				gids := strings.Split(order, "to")
				if len(gids) == 2 {
					fromGid, err := strconv.ParseInt(gids[0], 10, 64)
					if err != nil {
						msg.Text = "复制的群组ID有误"
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
					} else {
						Gid, _ := strconv.ParseInt(gids[1], 10, 64)
						rules := common.AllGroupRules[fromGid]
						db.UpdateGroupRule(Gid, rules.String())
						common.AddNewGroup(Gid)
						common.AllGroupRules[Gid] = rules
						msg.Text = "Copy all rules of the group to " + gids[1] + "from " + gids[0]
					}

				} else {
					fromGid, err := strconv.ParseInt(order, 10, 64)
					if err != nil {
						msg.Text = "复制的群组ID有误"
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
					} else {

						rules := common.AllGroupRules[fromGid]
						db.UpdateGroupRule(gid, rules.String())
						common.AddNewGroup(gid)
						common.AllGroupRules[gid] = rules
						msg.Text = "Copy all rules of the group to the current group " + "from " + order
					}
				}

			} else {
				msg.Text = copyText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessage(log, msg)
		} else {
			msg.Text = "该命令仅适用超级管理员"
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(log, msg)
		}
	case "del":
		if checkAdmin(log, gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				delRule(gid, order)
				msg.Text = "规则删除成功: " + order
			} else {
				msg.Text = delText
				msg.ParseMode = "Markdown"
			}
			sendMessage(log, msg)
		}
	case "delForAll":
		if checkSuperUser(log, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				for _, gid := range common.AllGroupId {
					delRule(gid, order)
				}
				msg.Text = "为所有群组或好友，删除规则成功: " + order
			} else {
				msg.Text = delForAllText
				msg.ParseMode = "Markdown"
			}
			sendMessage(log, msg)
		} else {
			msg.Text = "该命令仅适用超级管理员"
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(log, msg)
		}
	case "list":
		if checkAdmin(log, gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				if checkSuperUser(log, *upmsg.From) && upmsg.Chat.Type == "private" {

					order = strings.Replace(order, " ", "", -1)
					orderId, _ := strconv.ParseInt(order, 10, 64)
					rulelists := getRuleList(orderId)
					msg.Text = "ID: " + order
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true
					sendMessage(log, msg)
					if len(rulelists) > 0 {
						for _, rlist := range rulelists {
							msg.Text = rlist
							msg.ParseMode = "Markdown"
							msg.DisableWebPagePreview = true
							sendMessage(log, msg)
						}
					} else {
						msg.Text = "no rules."
						sendMessage(log, msg)
					}
				} else {
					msg.Text = "need SuperUser and private chat with bot."
					msg.ParseMode = "Markdown"
					msg.DisableWebPagePreview = true
					sendMessage(log, msg)
					deleteMessage(log, gid, upmsg.MessageID)
				}
			} else {
				rulelists := getRuleList(gid)
				msg.Text = "ID: " + strconv.FormatInt(gid, 10)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				sendMessage(log, msg)
				if len(rulelists) > 0 {
					for _, rlist := range rulelists {
						msg.Text = rlist
						msg.ParseMode = "Markdown"
						msg.DisableWebPagePreview = true
						sendMessage(log, msg)
					}
				} else {
					msg.Text = "no rules."
					sendMessage(log, msg)
				}
			}

		}
	case "admin":
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(uid, 10) + ") 请求管理员出来打屁股\r\n\r\n" + getAdmins(log, gid)
		msg.ParseMode = "Markdown"
		sendMessage(log, msg)
		banMember(log, gid, uid, 30)
	case "banme":
		botme, _ := bot.GetChatMember(api.GetChatMemberConfig{ChatConfigWithUser: api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID}})
		if botme.CanRestrictMembers {
			rand.Seed(time.Now().UnixNano())
			sec := rand.Intn(540) + 60
			banMember(log, gid, uid, int64(sec))
			msg.Text = "恭喜[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ")获得" + strconv.Itoa(sec) + "秒的禁言礼包"
			msg.ParseMode = "Markdown"
		} else {
			msg.Text = "请给我禁言权限,否则无法进行游戏"
		}
		sendMessage(log, msg)
	case "me":
		myuser := upmsg.From
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") 的账号信息" +
			"\r\nID: " + strconv.Itoa(int(uid)) +
			"\r\nUseName: [" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ")" +
			"\r\nLastName: " + myuser.LastName +
			"\r\nFirstName: " + myuser.FirstName +
			"\r\nIsBot: " + strconv.FormatBool(myuser.IsBot)
		msg.ParseMode = "Markdown"
		sendMessage(log, msg)
	default:
	}

}

// processReplyCommand 回复指令处理
func processReplyCommand(log logger.Logger, update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	// 回复类型的管理命令
	if upmsg.ReplyToMessage != nil {
		replyToUserId := upmsg.ReplyToMessage.From.ID
		switch upmsg.Text {
		case "kick":
			if checkAdmin(log, gid, *upmsg.From) {
				banMember(log, gid, replyToUserId, -1)
				mem, _ := bot.GetChatMember(api.GetChatMemberConfig{ChatConfigWithUser: api.ChatConfigWithUser{ChatID: gid, SuperGroupUsername: "", UserID: replyToUserId}})
				if !mem.CanSendMessages {
					msg = api.NewMessage(gid, "")
					msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") 将 " +
						"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.FormatInt(replyToUserId, 10) + ") " + "禁言了"
					msg.ParseMode = "Markdown"
					sendMessage(log, msg)
				}
			}
		case "unkick":
			if checkAdmin(log, gid, *upmsg.From) {
				unbanMember(log, gid, replyToUserId)
				// mem,_ := bot.GetChatMember(api.ChatConfigWithUser{gid, "", replyToUserId})
				//
				msg = api.NewMessage(gid, "")
				msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") 将 " +
					"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.FormatInt(replyToUserId, 10) + ") " + "解除禁言了"
				msg.ParseMode = "Markdown"
				sendMessage(log, msg)
			}

		//case "kick":
		//	if checkAdmin(log, gid, *upmsg.From) {
		//		kickMember(log, gid, replyToUserId)
		//		msg = api.NewMessage(gid, "")
		//		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") 关小黑屋了 " +
		//			"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.FormatInt(replyToUserId, 10) + ") "
		//		msg.ParseMode = "Markdown"
		//		sendMessage(log, msg)
		//	}
		//case "unkick":
		//	if checkAdmin(log, gid, *upmsg.From) {
		//		unkickMember(log, gid, replyToUserId)
		//		msg = api.NewMessage(gid, "")
		//		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.FormatInt(upmsg.From.ID, 10) + ") 放出来了 " +
		//			"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.FormatInt(replyToUserId, 10) + ") "
		//		msg.ParseMode = "Markdown"
		//		sendMessage(log, msg)
		//	}
		default:
		}
	}

}
