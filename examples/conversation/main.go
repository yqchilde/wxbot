package conversation

import (
	"fmt"
	"time"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

var selected = []string{"苹果", "香蕉", "火龙果"}

func init() {
	engine := control.Register("conversation", &control.Options{
		Alias: "会话",
		Help:  "可用于会话选择，比如注册时信息，或者选择信息等情况",
	})

	engine.OnFullMatch("选择水果").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		replyMsg := "请选择水果："
		for i := range selected {
			replyMsg += fmt.Sprintf("\n%d. %s", i+1, selected[i])
		}

		// ctx.CheckUserSession() 会话连续性，且必须同一个群同一个用户
		// ctx.CheckGroupSession() 会话连续性，且必须同一个群，可以不同用户，简单说就是一个群内多个用户在一个会话中
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText(replyMsg)
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				msg := ctx.MessageString()
				switch msg {
				case "1":
					ctx.ReplyTextAndAt("你选择了苹果")
					return
				case "2":
					ctx.ReplyTextAndAt("你选择了香蕉")
					return
				case "3":
					ctx.ReplyTextAndAt("你选择了火龙果")
					return
				}
			}
		}
	})
}
