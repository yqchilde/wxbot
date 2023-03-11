package robot

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var (
	bot         *Bot         // 当前机器人
	eventBuffer *EventBuffer // 事件缓冲区
)

type Bot struct {
	self      *Self
	config    *Config
	framework IFramework
}

// Run 运行并阻塞主线程，等待事件
func Run(c *Config, f IFramework) {
	if c.BufferLen == 0 {
		c.BufferLen = 4096
	}
	if c.Latency == 0 {
		c.Latency = time.Second
	}
	if c.MaxProcessTime == 0 {
		c.MaxProcessTime = time.Minute * 3
	}
	if c.ServerPort == 0 {
		c.ServerPort = 9528
	}

	bot = &Bot{config: c, framework: f}
	bot.self = &Self{bot: bot, User: &User{}}
	if c.connHookStatus {
		bot.self.Init()
		for i := range c.SuperUsers {
			if bot.self.friends.GetByWxId(c.SuperUsers[i]) == nil {
				log.Warnf("[robot] 您设置的管理员[%s]并不是您的好友，请修改config.yaml", c.SuperUsers[i])
			}
		}

		log.Printf("[robot] 共获取到%d个好友", bot.self.friends.Count())
		log.Printf("[robot] 共获取到%d个群组", bot.self.groups.Count())
		log.Printf("[robot] 共获取到%d个公众号", bot.self.mps.Count())
	}
	log.Printf("[robot] 机器人%s开始工作", c.BotNickname)

	eventBuffer = NewEventBuffer(bot.config.BufferLen)
	eventBuffer.Loop(bot.config.Latency, bot.config.MaxProcessTime, processEventAsync)
	runServer(c)
}

func processEventAsync(event *Event, framework IFramework, maxWait time.Duration) {
	ctx := &Ctx{
		Bot:       bot,
		State:     State{},
		Event:     event,
		framework: framework,
	}
	matcherLock.Lock()
	if hasMatcherListChanged {
		matcherListForRanging = make([]*Matcher, len(matcherList))
		copy(matcherListForRanging, matcherList)
		hasMatcherListChanged = false
	}
	matcherLock.Unlock()
	preProcessMessageEvent(ctx, event)
	go match(ctx, matcherListForRanging, maxWait)
}

// match 延迟 (1~100ms) 再处理事件
func match(ctx *Ctx, matchers []*Matcher, maxWait time.Duration) {
	goRule := func(rule Rule) <-chan bool {
		ch := make(chan bool, 1)
		go func() {
			defer func() {
				close(ch)
				if err := recover(); err != nil {
					log.Errorf("[robot]执行Rule时运行时发生错误: %v\n%v", err, string(debug.Stack()))
				}
			}()
			ch <- rule(ctx)
		}()
		return ch
	}
	goHandler := func(h Handler) <-chan struct{} {
		ch := make(chan struct{}, 1)
		go func() {
			defer func() {
				close(ch)
				if err := recover(); err != nil {
					log.Errorf("[robot]执行Handler时运行时发生错误: %v\n%v", err, string(debug.Stack()))
				}
			}()
			h(ctx)
			ch <- struct{}{}
		}()
		return ch
	}
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Millisecond)
	t := time.NewTimer(maxWait)
	defer t.Stop()
loop:
	for _, matcher := range matchers {
		for k := range ctx.State {
			delete(ctx.State, k)
		}
		m := matcher.copy()
		ctx.matcher = m

		// 处理前置条件
		if m.Engine != nil {
			for _, handler := range m.Engine.preHandler {
				c := goRule(handler)
				for {
					select {
					case ok := <-c:
						if !ok {
							if m.Break {
								break loop
							}
							continue loop
						}
					case <-t.C:
						if m.NoTimeout {
							t.Reset(maxWait)
							continue
						}
						break loop
					}
					break
				}
			}
		}
		// 处理rule
		for _, rule := range m.Rules {
			c := goRule(rule)
			for {
				select {
				case ok := <-c:
					if !ok {
						if m.Break {
							break loop
						}
						continue loop
					}
				case <-t.C:
					if m.NoTimeout {
						t.Reset(maxWait)
						continue
					}
					break loop
				}
				break
			}
		}
		// 处理中间条件
		if m.Engine != nil {
			for _, handler := range m.Engine.midHandler {
				c := goRule(handler)
				for {
					select {
					case ok := <-c:
						if !ok {
							if m.Break {
								break loop
							}
							continue loop
						}
					case <-t.C:
						if m.NoTimeout {
							t.Reset(maxWait)
							continue
						}
						break loop
					}
					break
				}
			}
		}
		// 处理handler
		if m.Handler != nil {
			c := goHandler(m.Handler)
			for {
				select {
				case <-c:
				case <-t.C:
					if m.NoTimeout {
						t.Reset(maxWait)
						continue
					}
					break loop
				}
				break
			}
		}
		if matcher.Temp {
			matcher.Delete()
		}
		// 处理后置条件
		if m.Engine != nil {
			for _, handler := range m.Engine.postHandler {
				c := goHandler(handler)
				for {
					select {
					case <-c:
					case <-t.C:
						if m.NoTimeout {
							t.Reset(maxWait)
							continue
						}
						break loop
					}
					break
				}
			}
		}
		if m.Block {
			break loop
		}
	}
}

// preProcessMessageEvent 预处理消息事件
func preProcessMessageEvent(ctx *Ctx, e *Event) {
	switch e.Type {
	case EventPrivateChat:
		if ctx.IsReference() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])引用消息 ==> %v", e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsText() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])文本消息 ==> %v", e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsImage() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s]图片消息 ==> %v", e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsVoice() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])语音消息", e.FromName, e.FromWxId))
		} else if ctx.IsShareCard() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])名片消息", e.FromName, e.FromWxId))
		} else if ctx.IsVideo() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])视频消息", e.FromName, e.FromWxId))
		} else if ctx.IsMemePictures() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])动态表情消息", e.FromName, e.FromWxId))
		} else if ctx.IsLocation() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])地理位置消息", e.FromName, e.FromWxId))
		} else if ctx.IsApp() {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])应用消息", e.FromName, e.FromWxId))
		} else {
			log.Println(fmt.Sprintf("[回调]收到私聊(%s[%s])未整理消息 ==> %v", e.FromName, e.FromWxId, e.Message.Content))
		}
	case EventGroupChat:
		if ctx.IsReference() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])引用消息 ==> %v", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsText() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])文本消息 ==> %v", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsImage() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])图片消息 ==> %v", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId, e.Message.Content))
		} else if ctx.IsVoice() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])语音消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else if ctx.IsShareCard() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])名片消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else if ctx.IsVideo() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])视频消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else if ctx.IsMemePictures() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])动态表情消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else if ctx.IsLocation() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])地理位置消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else if ctx.IsApp() {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])应用消息", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId))
		} else {
			log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])>用户(%s[%s])未整理消息 ==> %v", e.FromGroupName, e.FromGroup, e.FromName, e.FromWxId, e.Message.Content))
		}
	case EventMPChat:
		log.Println(fmt.Sprintf("[回调]收到订阅公众号(%s[%s])消息", e.FromName, e.FromWxId))
	case EventSelfMessage:
		log.Println(fmt.Sprintf("[回调]收到自己发送的消息 ==> %v", e.Message.Content))
	case EventFriendVerify:
		log.Println(fmt.Sprintf("[回调]收到好友验证消息, wxId:%s, nick:%s, content:%s", e.FriendVerifyMessage.WxId, e.FriendVerifyMessage.Nick, e.FriendVerifyMessage.Content))
	case EventTransfer:
		if len(e.TransferMessage.Memo) > 0 {
			log.Println(fmt.Sprintf("[回调]收到转账消息, wxId:%s, money:%s, memo:%s", e.TransferMessage.FromWxId, e.TransferMessage.Money, e.TransferMessage.Memo))
		} else {
			log.Println(fmt.Sprintf("[回调]收到转账消息, wxId:%s, money:%s", e.TransferMessage.FromWxId, e.TransferMessage.Money))
		}
	case EventMessageWithdraw:
		if e.WithdrawMessage.FromType == 1 {
			log.Println(fmt.Sprintf("[回调]收到撤回私聊(%s)消息", e.WithdrawMessage.FromWxId))
		} else if e.WithdrawMessage.FromType == 2 {
			log.Println(fmt.Sprintf("[回调]收到撤回群聊(%s[%s])消息", e.WithdrawMessage.FromGroup, e.WithdrawMessage.FromWxId))
		}
	case EventSystem:
		log.Println(fmt.Sprintf("[回调]收到系统消息 ==> %s", e.Message.Content))
	}
}

// GetCtx 获取当前系统中的CTX
func GetCtx() *Ctx {
	t := time.NewTimer(3 * time.Minute)
	for {
		select {
		case <-t.C:
			log.Fatal("[robot] 获取CTX超时")
		default:
			if bot == nil {
				time.Sleep(time.Second)
				continue
			}
			if bot.framework == nil {
				time.Sleep(time.Second)
				continue
			}
			t.Stop()
			return &Ctx{framework: bot.framework}
		}
	}
}

// GetBot 获取机器人本身
func GetBot() *Bot {
	return bot
}

// GetConfig 获取机器人配置
func (b *Bot) GetConfig() *Config {
	return b.config
}

// Friends 从缓存中获取好友列表
func (b *Bot) Friends() Friends {
	return b.self.friends
}

// Groups 从缓存中获取群列表
func (b *Bot) Groups() Groups {
	return b.self.groups
}

// MPs 从缓存中获取公众号列表
func (b *Bot) MPs() MPs {
	return b.self.mps
}

// Users 从缓存中获取所有用户列表
func (b *Bot) Users() []*User {
	var users []*User
	users = append(users, b.self.friends.AsUsers()...)
	users = append(users, b.self.groups.AsUsers()...)
	users = append(users, b.self.mps.AsUsers()...)
	return users
}

// GetSelf 获取Self对象，Self对象包含了对用户、群、公众号的包装
func (b *Bot) GetSelf() (*Self, error) {
	if b.self == nil {
		return nil, errors.New("bot self is nil")
	}
	return b.self, nil
}
