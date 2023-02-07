package robot

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var (
	WxBot       *Robot       // 当前机器人
	eventBuffer *EventBuffer // 事件缓冲区
)

type Robot struct {
	BotConfig        *Config
	Framework        IFramework
	FriendsList      []*FriendInfo
	GroupList        []*GroupInfo
	SubscriptionList []*SubscriptionInfo
}

type Config struct {
	BotWxId        string        // 机器人微信ID
	BotNickname    string        // 机器人名称
	SuperUsers     []string      // 超级用户(管理员)
	CommandPrefix  string        // 管理员触发命令
	BufferLen      uint          // 事件缓冲区长度, 默认4096
	Latency        time.Duration // 事件处理延迟 (延迟 latency + (0~100ms) 再处理事件) (默认1s)
	MaxProcessTime time.Duration // 事件最大处理时间 (默认3min)
	Framework      IFramework    // 接入框架需实现该接口
}

// Init 初始化机器人
func Init(c *Config) *Robot {
	if c.BufferLen == 0 {
		c.BufferLen = 4096
	}
	if c.Latency == 0 {
		c.Latency = time.Second
	}
	if c.MaxProcessTime == 0 {
		c.MaxProcessTime = time.Minute * 3
	}
	go monitoringWechatData()

	return &Robot{
		BotConfig: c,
		Framework: c.Framework,
	}
}

// Run 运行并阻塞主线程，等待事件
func (r *Robot) Run() {
	eventBuffer = NewEventBuffer(r.BotConfig.BufferLen)
	eventBuffer.Loop(r.BotConfig.Latency, r.BotConfig.MaxProcessTime, processEventAsync)
	r.Framework.Callback(eventBuffer.ProcessEvent)
}

func processEventAsync(event *Event, framework IFramework, maxWait time.Duration) {
	ctx := &Ctx{
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
	preProcessMessageEvent(event)
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
func preProcessMessageEvent(e *Event) {
	switch e.Type {
	case EventPrivateChat:
		log.Println(fmt.Sprintf("[回调]收到私聊(%s)消息 ==> %v", e.FromName, e.Message.Content))
	case EventGroupChat:
		log.Println(fmt.Sprintf("[回调]收到群聊(%s[%s])消息 ==> %v", e.FromGroupName, e.FromWxId, e.Message.Content))
	case EventSubscription:
		log.Println(fmt.Sprintf("[回调]收到订阅公众号(%s)消息", e.FromWxId))
	case EventSelfMessage:
		log.Println(fmt.Sprintf("[回调]收到自己发送的消息 ==> %v", e.Message.Content))
	case EventFriendVerify:
		log.Println(fmt.Sprintf("[回调]收到好友验证消息, wxId:%s, nick:%s, content:%s", e.FriendVerify.WxId, e.FriendVerify.Nick, e.FriendVerify.Content))
	case EventTransfer:
		if len(e.Transfer.Memo) > 0 {
			log.Println(fmt.Sprintf("[回调]收到转账消息, wxId:%s, money:%s, memo:%s", e.Transfer.FromWxId, e.Transfer.Money, e.Transfer.Memo))
		} else {
			log.Println(fmt.Sprintf("[回调]收到转账消息, wxId:%s, money:%s", e.Transfer.FromWxId, e.Transfer.Money))
		}
	case EventMessageWithdraw:
		if e.Withdraw.FromType == 1 {
			log.Println(fmt.Sprintf("[回调]收到撤回私聊(%s)消息 ==> %s", e.Withdraw.FromWxId, e.Withdraw.Msg))
		} else if e.Withdraw.FromType == 2 {
			log.Println(fmt.Sprintf("[回调]收到撤回群聊(%s[%s])消息 ==> %s", e.Withdraw.FromGroup, e.Withdraw.FromWxId, e.Withdraw.Msg))
		}
	case EventSystem:
		log.Println(fmt.Sprintf("[回调]收到系统消息 ==> %s", e.Message.Content))
	}
}

// monitoringWechatData 监控微信数据
func monitoringWechatData() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		friendsList, _ := WxBot.Framework.GetFriendsList(true)
		groupList, _ := WxBot.Framework.GetGroupList(true)
		subscriptionList, _ := WxBot.Framework.GetSubscriptionList(true)
		WxBot.FriendsList = friendsList
		WxBot.GroupList = groupList
		WxBot.SubscriptionList = subscriptionList
	}
}

// GetCTX 获取当前系统中的CTX
func GetCTX() *Ctx {
	t := time.NewTimer(3 * time.Minute)
	for {
		select {
		case <-t.C:
			log.Fatal("[robot] 获取CTX超时")
		default:
			if WxBot == nil {
				time.Sleep(time.Second)
				continue
			}
			if WxBot.Framework == nil {
				time.Sleep(time.Second)
				continue
			}
			t.Stop()
			return &Ctx{framework: WxBot.Framework}
		}
	}
}
