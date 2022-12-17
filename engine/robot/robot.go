package robot

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var (
	eveRing   eventRing // 用于存储事件的环形队列
	BotConfig *Config   // 机器人配置
)

type Config struct {
	BotWxId        string        // 机器人微信ID
	BotNickname    string        // 机器人名称
	SuperUsers     []string      // 超级用户
	CommandPrefix  string        // 触发命令
	RingLen        uint          // 事件环长度 (默认4096)
	Latency        time.Duration // 事件处理延迟 (延迟 latency + (0~100ms) 再处理事件) (默认1min)
	MaxProcessTime time.Duration // 事件最大处理时间 (默认3min)
	Framework      IFramework    // 接入框架需实现该接口
}

// Run 主函数，启动机器人
func Run(c *Config) {
	if c.RingLen == 0 {
		c.RingLen = 4096
	}
	if c.Latency == 0 {
		c.Latency = time.Second
	}
	if c.MaxProcessTime == 0 {
		c.MaxProcessTime = time.Minute * 3
	}
	BotConfig = c
	eveRing = newRing(c.RingLen)
	eveRing.loop(c.Latency, c.MaxProcessTime, processEventAsync)
	log.Printf("[robot] 机器人%s开始工作", c.BotNickname)
	c.Framework.Callback(eveRing.processEvent)
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
						log.Debug("[robot] preHandler处理达到最大时延, 退出")
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
					log.Debug("[robot] rule处理达到最大时延, 退出")
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
						log.Debug("[robot] midHandler处理达到最大时延, 退出")
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
					log.Debug("[robot] Handler处理达到最大时延, 退出")
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
						log.Warn("[robot] postHandler处理达到最大时延, 退出")
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
		log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", e.FromWxId, e.Message.Content))
	case EventGroupChat:
		log.Println(fmt.Sprintf("收到群聊(%s[%s])消息 ==> %v", e.FromGroup, e.FromWxId, e.Message.Content))
	case EventFriendVerify:
		log.Println(fmt.Sprintf("收到好友验证消息, wxId:%s, nick:%s, content:%s", e.FriendVerify.WxId, e.FriendVerify.Nick, e.FriendVerify.Content))
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
			if BotConfig != nil {
				t.Stop()
				return &Ctx{framework: BotConfig.Framework}
			}
			time.Sleep(time.Second)
		}
	}
}
