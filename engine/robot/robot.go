package robot

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/yqchilde/pkgs/log"
)

func init() {
	log.Default(2)
}

var BotConfig Config

// Config 关于机器人的相关参数配置
type Config struct {
	Nickname       string        // 机器人名称
	SuperUsers     []string      // 超级用户
	CommandPrefix  string        // 触发命令
	RingLen        uint          // 事件环长度 (默认4096)
	Latency        time.Duration // 事件处理延迟 (延迟 latency + (0~100ms) 再处理事件) (默认1min)
	MaxProcessTime time.Duration // 事件最大处理时间 (默认3min)
	Framework      Framework     // 接入框架
}

// Framework 机器人通信驱动，暂时只支持HTTP
type Framework interface {
	Callback(func(*Event, APICaller))
}

// APICaller 定义了机器人的API调用接口，接入的框架需要实现这个接口
type APICaller interface {
	// SendText 发送文本消息
	// toWxId: 好友ID/群ID
	// text: 文本内容
	SendText(toWxId, text string) error

	// SendTextAndAt 发送文本消息并@，只有群聊有效
	// toGroupWxId: 群ID
	// toWxId: 好友ID/群ID/all
	// toWxName: 好友昵称/群昵称，留空为自动获取
	// text: 文本内容
	SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error

	// SendImage 发送图片消息
	// toWxId: 好友ID/群ID
	// path: 图片路径
	SendImage(toWxId, path string) error

	// SendShareLink 发送分享链接消息
	// toWxId: 好友ID/群ID
	// title: 标题
	// desc: 描述
	// imageUrl: 图片链接
	// jumpUrl: 跳转链接
	SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error
}

// 事件环
var eveRing eventRing

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
	BotConfig = *c
	eveRing = newRing(c.RingLen)
	eveRing.loop(c.Latency, c.MaxProcessTime, processEventAsync)
	log.Printf("[robot] 机器人%s开始工作", c.Nickname)
	c.Framework.Callback(eveRing.processEvent)
}

func processEventAsync(event *Event, caller APICaller, maxWait time.Duration) {
	ctx := &Ctx{
		State:  State{},
		Event:  event,
		caller: caller,
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
			// clear state
			delete(ctx.State, k)
		}
		m := matcher.copy()
		ctx.matcher = m

		// pre handler
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
						log.Warn("[robot] preHandler处理达到最大时延, 退出")
						break loop
					}
					break
				}
			}
		}

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
					log.Warn("[robot] rule处理达到最大时延, 退出")
					break loop
				}
				break
			}
		}

		// mid handler
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
						log.Warn("[robot] midHandler处理达到最大时延, 退出")
						break loop
					}
					break
				}
			}
		}

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
					log.Warn("[robot] Handler处理达到最大时延, 退出")
					break loop
				}
				break
			}
		}
		if matcher.Temp {
			matcher.Delete()
		}

		if m.Engine != nil {
			// post handler
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
func preProcessMessageEvent(ctx *Ctx, e *Event) {
	if ctx.IsSendByPrivateChat() {
		log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", e.Message.FromWxId, e.Message.Msg))
	} else if ctx.IsSendByGroupChat() {
		log.Println(fmt.Sprintf("收到群聊(%s[%s])消息 ==> %v", e.Message.FromGroup, e.Message.FromWxId, e.Message.Msg))
	}
}
