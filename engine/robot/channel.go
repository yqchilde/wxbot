package robot

// EventChannel 用于异步获取指定事件
type EventChannel struct {
	Type     string
	Priority uint64
	Rule     []Rule
	Block    bool
}

// NewEventChannel 创建一个新的 EventChannel用于异步获取指定事件
func NewEventChannel(Priority uint64, Block bool, rule ...Rule) *EventChannel {
	return &EventChannel{
		Priority: Priority,
		Rule:     rule,
		Block:    Block,
	}
}

// EventChannel 用当前上下文创建一个EventChannel，会阻塞其他事件
func (m *Matcher) EventChannel(rule ...Rule) *EventChannel {
	return &EventChannel{
		Priority: m.Priority,
		Block:    m.Block,
		Rule:     rule,
	}
}

// Next 返回一个 chan 用于接收下一个指定事件
// 该 chan 必须接收，如需手动取消监听，请使用 Repeat 方法
func (n *EventChannel) Next() <-chan *Ctx {
	ch := make(chan *Ctx, 1)
	StoreTempMatcher(&Matcher{
		Block:    n.Block,
		Priority: n.Priority,
		Rules:    n.Rule,
		Engine:   defaultEngine,
		Handler: func(ctx *Ctx) {
			ch <- ctx
			close(ch)
		},
	})
	return ch
}

// Repeat 返回一个 chan 用于接收无穷个指定事件，和一个取消监听的函数
// 如果没有取消监听，将不断监听指定事件
func (n *EventChannel) Repeat() (recv <-chan *Ctx, cancel func()) {
	ch, done := make(chan *Ctx, 1), make(chan struct{})
	go func() {
		in := make(chan *Ctx, 1)
		matcher := StoreMatcher(&Matcher{
			Block:    n.Block,
			Priority: n.Priority,
			Rules:    n.Rule,
			Engine:   defaultEngine,
			Handler: func(ctx *Ctx) {
				in <- ctx
			},
		})
		for {
			select {
			case e := <-in:
				ch <- e
			case <-done:
				matcher.Delete()
				close(in)
				close(ch)
				return
			}
		}
	}()
	return ch, func() {
		done <- struct{}{}
	}
}
