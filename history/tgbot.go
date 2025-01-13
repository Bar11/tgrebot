package history

import (
	"github.com/panjf2000/ants/v2"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/chain5j/logger"
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"tg-keyword-reply-bot/common"
	"tg-keyword-reply-bot/db"
)

var (
	chPool *ants.Pool
)

// start 开始工作
func start(log logger.Logger, botToken string) {
	var err error

	bot, err := api.NewBotAPI(botToken)
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

	updates, err := bot.GetUpdatesChan(u)
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
	log.Debug("update msg", "msg", upmsg.Text)
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	log.Debug("message from update", "gid", gid, "uid", uid, "mid", upmsg.MessageID, "uname", upmsg.From.String())
	// 检查是不是新加的群或者新开的人
	in := checkInGroup(gid)
	if !in {
		// 不在就需要加入, 内存中加一份, 数据库中添加一条空规则记录
		common.AddNewGroup(gid)
		db.AddNewGroup(gid)
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
			for _, auser := range *(upmsg.NewChatMembers) {
				if checkQingzhen(auser.UserName) ||
					checkQingzhen(auser.FirstName) ||
					checkQingzhen(auser.LastName) {
					banMember(log, gid, uid, -1)
				}
			}
		}
		// 检查清真并剔除
		if checkQingzhen(upmsg.Text) {
			_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
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
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
	} else if strings.HasPrefix(replyText, "ban") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		banMember(log, gid, uid, -1)
	} else if strings.HasPrefix(replyText, "kick") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		kickMember(log, gid, uid)
	} else if strings.HasPrefix(replyText, "photo:") {
		sendPhoto(log, gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "gif:") {
		sendGif(log, gid, replyText[4:])
	} else if strings.HasPrefix(replyText, "video:") {
		sendVideo(log, gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "file:") {
		sendFile(log, gid, replyText[5:])
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
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(update.Message.Chat.ID, upmsg.MessageID))
	log.Info("bot delete the msg", "gid", gid, "uid", uid, "mid", upmsg.MessageID)
	switch upmsg.Command() {
	case "start", "help":
		msg.Text = "本机器人能够自动回复特定关键词"
		sendMessage(log, msg)
	case "add":
		if checkAdmin(log, gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				addRule(gid, order)
				msg.Text = "规则添加成功: " + order
			} else {
				msg.Text = addText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
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
	case "list":
		if checkAdmin(log, gid, *upmsg.From) {
			rulelists := getRuleList(gid)
			msg.Text = "ID: " + strconv.FormatInt(gid, 10)
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(log, msg)
			for _, rlist := range rulelists {
				msg.Text = rlist
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				sendMessage(log, msg)
			}
		}
	case "admin":
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(uid) + ") 请求管理员出来打屁股\r\n\r\n" + getAdmins(log, gid)
		msg.ParseMode = "Markdown"
		sendMessage(log, msg)
		banMember(log, gid, uid, 30)
	case "banme":
		botme, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID})
		if botme.CanRestrictMembers {
			rand.Seed(time.Now().UnixNano())
			sec := rand.Intn(540) + 60
			banMember(log, gid, uid, int64(sec))
			msg.Text = "恭喜[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")获得" + strconv.Itoa(sec) + "秒的禁言礼包"
			msg.ParseMode = "Markdown"
		} else {
			msg.Text = "请给我禁言权限,否则无法进行游戏"
		}
		sendMessage(log, msg)
	case "me":
		myuser := upmsg.From
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 的账号信息" +
			"\r\nID: " + strconv.Itoa(uid) +
			"\r\nUseName: [" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")" +
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
		case "ban":
			if checkAdmin(log, gid, *upmsg.From) {
				banMember(log, gid, replyToUserId, -1)
				mem, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, SuperGroupUsername: "", UserID: replyToUserId})
				if !mem.CanSendMessages {
					msg = api.NewMessage(gid, "")
					msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 禁言了 " +
						"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(replyToUserId) + ") "
					msg.ParseMode = "Markdown"
					sendMessage(log, msg)
				}
			}
		case "unban":
			if checkAdmin(log, gid, *upmsg.From) {
				unbanMember(log, gid, replyToUserId)
				// mem,_ := bot.GetChatMember(api.ChatConfigWithUser{gid, "", replyToUserId})
				//
				msg = api.NewMessage(gid, "")
				msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 解禁了 " +
					"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(replyToUserId) + ") "
				msg.ParseMode = "Markdown"
				sendMessage(log, msg)
			}
		case "kick":
			if checkAdmin(log, gid, *upmsg.From) {
				kickMember(log, gid, replyToUserId)
			}
		case "unkick":
			if checkAdmin(log, gid, *upmsg.From) {
				unkickMember(log, gid, replyToUserId)
			}
		default:
		}
	}
}
