package robot

import (
	"sync"
)

type Ctx struct {
	matcher   *Matcher
	Event     *Event
	State     State
	framework IFramework

	// lazy message
	once    sync.Once
	message string
}

// GetMatcher 获取匹配器
func (ctx *Ctx) GetMatcher() *Matcher {
	return ctx.matcher
}

// MessageString 字符串消息便于Regex
func (ctx *Ctx) MessageString() string {
	ctx.once.Do(func() {
		if ctx.Event != nil {
			ctx.message = ctx.Event.Message.Msg
		}
	})
	return ctx.message
}

// CheckSession 判断会话连续性
func (ctx *Ctx) CheckSession() Rule {
	return func(ctx2 *Ctx) bool {
		return ctx.Event.FromWxId == ctx2.Event.FromWxId &&
			ctx.Event.FromGroup == ctx2.Event.FromGroup
	}
}

// EventChannel 用当前事件创建一个新的事件通道
func (ctx *Ctx) EventChannel(rule ...Rule) *EventChannel {
	return ctx.matcher.EventChannel(rule...)
}
