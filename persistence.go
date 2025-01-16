package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"tg-keyword-reply-bot/common"
	"tg-keyword-reply-bot/db"
)

const (
	helpText = "管理员命令:\r\n" +
		"/add\r\n" + addText +
		"/del\r\n" + delText +
		"/list\r\n" + listText +
		"/admin\r\n" + "开发中。。。\r\n\r\n" +
		"/me\r\n" + "查看本人账号信息。。。\r\n\r\n" +
		"超级管理员命令:\r\n" +
		"/addForAll\r\n" + addForAllText +
		"/delForAll\r\n" + delForAllText +
		"/copy\r\n" + copyText

	addText = "格式要求:\r\n" +
		"`/add 关键字===回复内容`\r\n" +
		"`/add 关键字1||关键字2===回复内容`\r\n" +
		"例如:\r\n" +
		"`/add 机场===https://jiji.cool`\r\n" +
		"就会添加一条规则, 关键词是机场, 回复内容是网址\r\n\r\n"
	addForAllText = "权限要求:\r\n" +
		"超级管理员\r\n" +
		"格式要求:\r\n" +
		"`/addForAll 关键字===回复内容`\r\n" +
		"`/addForAll 关键字1||关键字2===回复内容`\r\n" +
		"例如:\r\n" +
		"`/addForAll 机场===https://jiji.cool`\r\n" +
		"就会为所有群组和好友添加一条规则, 关键词是机场, 回复内容是网址\r\n\r\n"
	listText = "可以查看本群`group-id`和所有自动回复规则\r\n" +
		"例如:\r\n" +
		"`/list `\r\n" +
		"ID:0123456789\r\n" +
		"回复规则。。。。。。\r\n\r\n"
	delText = "格式要求:\r\n" +
		"`/del 关键字`\r\n" +
		"例如:\r\n" +
		"`/del 机场`\r\n" +
		"就会删除一条规则,机器人不再回复机场关键词\r\n\r\n"
	delForAllText = "权限要求:\r\n" +
		"超级管理员\r\n" +
		"格式要求:\r\n" +
		"`/delForAll 关键字`\r\n" +
		"例如:\r\n" +
		"`/delForAll 机场`\r\n" +
		"就会删除所有群组和好友下的这一条规则,机器人所有群组和好友不再回复机场关键词\r\n\r\n"
	copyText = "权限要求:\r\n" +
		"超级管理员\r\n" +
		"获取group-id的方法：\r\n" +
		"在需要复制的群组或好友窗口下输入/list命令\r\n" +
		"格式要求:\r\n" +
		"`/copy grup-id`\r\n" +
		"例如:\r\n" +
		"`/copy 1234567890`\r\n" +
		"就会复制ID为：1234567890的群组或好友下的所有规则，到当前群组\r\n\r\n"
)

// addRule 添加规则
func addRule(gid int64, rule string) {
	rules := common.AllGroupRules[gid]
	r := strings.Split(rule, "===")
	if len(r) < 2 {
		return
	}
	keys, value := r[0], r[1]
	if strings.Contains(keys, "||") {
		ks := strings.Split(keys, "||")
		for _, key := range ks {
			_addOneRule(key, value, rules)
		}
	} else {
		_addOneRule(keys, value, rules)
	}
	db.UpdateGroupRule(gid, rules.String())
}

// 给addRule使用的辅助方法
func _addOneRule(key string, value string, rules common.RuleMap) {
	key = strings.Replace(key, " ", "", -1)
	rules[key] = value
}

// 删除规则
func delRule(gid int64, key string) {
	rules := common.AllGroupRules[gid]
	delete(rules, key)
	db.UpdateGroupRule(gid, rules.String())
}

// 获取一个群组所有规则的列表
func getRuleList(gid int64) []string {
	kvs := common.AllGroupRules[gid]
	str := ""
	var strs []string
	num := 1
	group := 0
	for k, v := range kvs {
		str += "\r\n\r\n规则" + strconv.Itoa(num) + ":\r\n`" + k + "` => `" + v + "` "
		num++
		group++
		if group == 10 {
			group = 0
			strs = append(strs, str)
			str = ""
		}
	}
	strs = append(strs, str)
	return strs
}

// 查询是否包含相应的自动回复规则
func findKey(gid int64, input string) string {
	kvs := common.AllGroupRules[gid]
	fmt.Println("AllGroupRules:", kvs)
	for keyword, reply := range kvs {
		if strings.HasPrefix(keyword, "re:") {
			keyword = keyword[3:]
			match, _ := regexp.MatchString(keyword, input)
			if match {
				return reply
			}
		} else if strings.Contains(input, keyword) {
			return reply
		}
	}
	return ""
}
