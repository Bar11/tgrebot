package db

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-keyword-reply-bot/common"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // 初始化gorm使用sqlite
)

var db *gorm.DB

type setting struct {
	gorm.Model
	Key   string `gorm:"unique;not null"`
	Value string
}

type rule struct {
	gorm.Model
	GroupId  int64 `gorm:"unique;not null"`
	RuleJson string
}

type messageRecord struct {
	gorm.Model
	MessageID     int    `json:"message_id"`
	From          string `json:"from_username"`
	FromId        int64  `json:"from_id"`
	FromLang      string `json:"from_lang"`
	Date          int    `json:"date"`
	ChatId        int64  `json:"chat_id"`
	CharType      string `json:"char_type"`
	ChatTitle     string `json:"chat_title"`
	ForwardFromId int64  `json:"forward_from_id"`
	ForwardFrom   string `json:"forward_from_userName"`
	ForwardDate   int    `json:"forward_date"`
	ForwardLang   string `json:"forward_lang"`
	Text          string `json:"text"`
	Photo         bool   `json:"photo"`
	Document      bool   `json:"document"`
	Vedio         bool   `json:"vedio"`
	Voice         bool   `json:"voice"`
	ReplyId       int    `json:"reply_id"`
}

// Init 数据库初始化，包括新建数据库（如果还没有建立），基本数据的读写
func Init(newToken string) (token string) {
	dbtmp, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to connect database")
	}
	db = dbtmp
	db.AutoMigrate(&setting{}, &rule{}, &messageRecord{})
	var tokenSetting setting
	db.Find(&tokenSetting, "Key=?", "token")
	token = tokenSetting.Value
	if newToken != "" {
		token = newToken
		if tokenSetting.ID > 0 {
			tokenSetting.Value = newToken
			db.Model(&tokenSetting).Update(tokenSetting)
		} else {
			db.Create(&setting{
				Key:   "token",
				Value: newToken,
			})
		}
	}
	readAllGroupRules()
	return
}

func AddMessageRecord(message api.Message) {
	var messageForwardFromId int64 = 0
	var messageForwardDate = 0
	var messageForwardLang = ""
	var messageForwardFrom = ""
	var replyToMessage = 0
	if message.ForwardFrom != nil {
		messageForwardFromId = message.ForwardFrom.ID
		messageForwardDate = message.ForwardDate
		messageForwardLang = message.ForwardFrom.LanguageCode
		messageForwardFrom = message.ForwardFrom.UserName
	}
	if message.ReplyToMessage != nil {
		replyToMessage = message.ReplyToMessage.MessageID
	}

	db.Create(&messageRecord{
		MessageID:     message.MessageID,
		From:          message.From.FirstName + message.From.LastName,
		FromId:        message.From.ID,
		Date:          message.Date,
		ChatId:        message.Chat.ID,
		ChatTitle:     message.Chat.Title,
		Text:          message.Text,
		FromLang:      message.From.LanguageCode,
		CharType:      message.Chat.Type,
		ForwardFromId: messageForwardFromId,
		ForwardFrom:   messageForwardFrom,
		ForwardDate:   messageForwardDate,
		ForwardLang:   messageForwardLang,
		Photo:         message.Photo != nil,
		Document:      message.Document != nil,
		Vedio:         message.Video != nil,
		Voice:         message.Voice != nil,
		ReplyId:       replyToMessage,
	})
}

// AddNewGroup 数据库中添加一条记录来记录新群组的规则
func AddNewGroup(groupId int64) {
	db.Create(&rule{
		GroupId:  groupId,
		RuleJson: "",
	})
}

// UpdateGroupRule 更新群组的规则
func UpdateGroupRule(groupId int64, ruleJson string) {
	db.Model(&rule{}).Where("group_id=?", groupId).Update("rule_json", ruleJson)
}

// 读取所有的规则到内容中
func readAllGroupRules() {
	var allGroupRules []rule
	db.Find(&allGroupRules)
	for _, rule := range allGroupRules {
		ruleStruct := common.Json2kvs(rule.RuleJson)
		common.AllGroupRules[rule.GroupId] = ruleStruct
		common.AllGroupId = append(common.AllGroupId, rule.GroupId)
	}
}
