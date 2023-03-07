package friendadd

import (
	"fmt"
	"strings"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("friendadd", &control.Options{
		Alias: "自动通过好友添加请求",
	})

	engine.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
		// 监听加好友事件
		if ctx.IsEventFriendVerify() {
			f := ctx.Event.FriendVerify
			// 判断一下好友验证消息
			nickname := robot.GetBot().GetConfig().BotNickname
			if strings.ToLower(f.Content) != strings.ToLower(nickname) {
				return
			}
			if err := ctx.AgreeFriendVerify(f.V3, f.V4, f.Scene); err != nil {
				log.Errorf("同意好友请求失败: %v", err)
				return
			}
			ctx.SendText(f.WxId, fmt.Sprintf("你好，我是%s，感谢你发现并使用我，输入\"开始会话\"体验 ChatGPT！", nickname))
		}
	})
}
