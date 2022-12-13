package robot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yqchilde/pkgs/log"
)

// PrefixRule 检查消息前缀
func PrefixRule(prefixes ...string) Rule {
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		msg := ctx.MessageString()
		for _, prefix := range prefixes {
			if strings.HasPrefix(msg, prefix) {
				ctx.State["prefix"] = prefix
				arg := strings.TrimLeft(msg[len(prefix):], " ")
				ctx.State["args"] = arg
				return true
			}
		}
		return false
	}
}

// SuffixRule 检查消息后缀
func SuffixRule(suffixes ...string) Rule {
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		msg := ctx.MessageString()
		for _, suffix := range suffixes {
			if strings.HasSuffix(msg, suffix) {
				ctx.State["suffix"] = suffix
				arg := strings.TrimRight(msg[:len(msg)-len(suffix)], " ")
				ctx.State["args"] = arg
				return true
			}
		}
		return false
	}
}

// CommandRule 检查消息是否为命令
func CommandRule(commands ...string) Rule {
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		if !strings.HasPrefix(ctx.Event.Message.Msg, BotConfig.CommandPrefix) {
			return false
		}
		cmdMessage := ctx.Event.Message.Msg[len(BotConfig.CommandPrefix):]
		for _, command := range commands {
			if strings.HasPrefix(cmdMessage, command) {
				ctx.State["command"] = command
				arg := strings.TrimLeft(cmdMessage[len(command):], " ")
				ctx.State["args"] = arg
				return true
			}
		}
		return false
	}
}

// RegexRule 检查消息是否匹配正则表达式
func RegexRule(regexPattern string) Rule {
	regex := regexp.MustCompile(regexPattern)
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		msg := ctx.MessageString()
		if matched := regex.FindStringSubmatch(msg); matched != nil {
			ctx.State["regex_matched"] = matched
			return true
		}
		return false
	}
}

// KeywordRule 检查消息是否包含关键字
func KeywordRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		msg := ctx.MessageString()
		for _, str := range src {
			if strings.Contains(msg, str) {
				ctx.State["keyword"] = str
				return true
			}
		}
		return false
	}
}

// FullMatchRule 检查消息是否完全匹配
func FullMatchRule(src ...string) Rule {
	return func(ctx *Ctx) bool {
		if !ctx.IsText() {
			return false
		}
		msg := ctx.MessageString()
		for _, str := range src {
			if str == msg {
				ctx.State["matched"] = msg
				return true
			}
		}
		return false
	}
}

// AdminPermission 只允许系统配置的管理员使用
func AdminPermission(ctx *Ctx) bool {
	for _, su := range BotConfig.SuperUsers {
		if su == ctx.Event.FromWxId {
			return true
		}
	}
	return false
}

// UserOrGroupAdmin 允许用户单独使用或群管使用
func UserOrGroupAdmin(ctx *Ctx) bool {
	if ctx.IsSendByGroupChat() {
		return AdminPermission(ctx)
	}
	return true
}

// HasMemePicture 检查消息是否存在表情包图片
func HasMemePicture(ctx *Ctx) bool {
	url, has := ctx.IsMemePictures()
	if has {
		ctx.State["image_url"] = url
		return true
	}
	return false
}

// MustMemePicture 消息不存在表情包图片阻塞至有图片，默认30s，超时返回false
// 阻塞时长可通过ctx.State["timeout"]设置
func MustMemePicture(ctx *Ctx) bool {
	if HasMemePicture(ctx) {
		return true
	}
	var timeout time.Duration
	if t, ok := ctx.State["timeout"]; ok {
		if v, ok := t.(time.Duration); !ok {
			log.Errorf("ctx.State[\"timeout\"] must be time.Duration")
			return false
		} else {
			timeout = v
		}
	} else {
		timeout = 30 * time.Second
	}
	ctx.ReplyTextAndAt(fmt.Sprintf("请在%d秒内发送表情包图片", int(timeout.Seconds())))
	next := NewEventChannel(999, true, ctx.CheckSession(), HasMemePicture).Next()
	select {
	case <-time.After(timeout):
		return false
	case newCtx := <-next:
		ctx.State["image_url"] = newCtx.State["image_url"]
		return true
	}
}

// OnlyGroup 只允许群聊使用
func OnlyGroup(ctx *Ctx) bool {
	return ctx.IsSendByGroupChat()
}

// OnlyPrivate 只允许私聊使用
func OnlyPrivate(ctx *Ctx) bool {
	return ctx.IsSendByPrivateChat()
}

// OnlyAtMe 只允许@机器人使用
func OnlyAtMe(ctx *Ctx) bool {
	return ctx.IsSendByGroupChat() && ctx.IsAt()
}
