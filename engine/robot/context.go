package robot

import (
	"regexp"
	"strings"
	"sync"
)

type Ctx struct {
	matcher   *Matcher
	Bot       *Bot
	Event     *Event
	State     State
	framework IFramework

	// lazy message
	once    sync.Once
	mutex   sync.Mutex
	message string
}

// GetMatcher 获取匹配器
func (ctx *Ctx) GetMatcher() *Matcher {
	return ctx.matcher
}

// MessageString 字符串消息便于Regex
func (ctx *Ctx) MessageString() string {
	ctx.once.Do(func() {
		if ctx.Event != nil && ctx.IsText() {
			if !ctx.IsAt() || ctx.IsEventPrivateChat() {
				ctx.message = ctx.Event.Message.Content
			} else {
				switch bot.config.Framework.Name {
				case "千寻", "qianxun", "Dean":
					ctx.message = strings.TrimPrefix(ctx.Event.Message.Content, "@"+bot.self.Nick)
					ctx.message = strings.TrimSpace(ctx.message)
				case "VLW", "vlw":
					regex := regexp.MustCompile(`\[at=.*\]\s*`)
					ctx.message = regex.ReplaceAllString(ctx.Event.Message.Content, "")
				}
			}
		}
	})
	return ctx.message
}

// CheckUserSession 判断会话连续性，必须同一个群同一个用户
func (ctx *Ctx) CheckUserSession() Rule {
	return func(ctx2 *Ctx) bool {
		return ctx.Event.FromWxId == ctx2.Event.FromWxId &&
			ctx.Event.FromGroup == ctx2.Event.FromGroup
	}
}

// CheckGroupSession 判断会话连续性，必须同一个群，可以不同用户
func (ctx *Ctx) CheckGroupSession() Rule {
	return func(ctx2 *Ctx) bool {
		return ctx.Event.FromUniqueID == ctx2.Event.FromUniqueID
	}
}

// EventChannel 用当前事件创建一个新的事件通道
func (ctx *Ctx) EventChannel(rule ...Rule) *EventChannel {
	return ctx.matcher.EventChannel(rule...)
}
