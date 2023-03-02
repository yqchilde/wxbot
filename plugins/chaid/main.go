package chaid

import (
	"fmt"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("chaid", &control.Options{
		Alias: "查ID",
		Help: "描述:\n" +
			"模糊查用户ID，仅支持查询已添加好友、群友、公众号\n\n" +
			"指令:\n" +
			"* 查id [昵称/备注]",
	})

	// 查系统ID，只能由管理员使用
	engine.OnRegex(`^查(?i:id) (.+)$`, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		id, reply := ctx.State["regex_matched"].([]string)[1], ""
		findUsers := ctx.FuzzyGetByRemarkOrNick(id)
		for _, user := range findUsers {
			if user.IsFriend() {
				reply += fmt.Sprintf("类别：%s\nID：%s\n昵称：%s\n\n", "好友", user.WxId, user.Nick)
			} else if user.IsGroup() {
				reply += fmt.Sprintf("类别：%s\nID：%s\n昵称：%s\n\n", "群组", user.WxId, user.Nick)
			} else if user.IsMP() {
				reply += fmt.Sprintf("类别：%s\nID：%s\n昵称：%s\n\n", "公众号", user.WxId, user.Nick)
			}
		}
		if reply == "" {
			reply = "未查到"
		}
		ctx.ReplyTextAndAt(reply)
	})
}
