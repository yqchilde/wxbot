package robot

// FutureEvent 交互的核心，用于异步获取指定事件
type FutureEvent struct {
	Type     string
	Priority int
	Rule     []Rule
	Block    bool
}

func NewFutureEvent(Priority int, Block bool, rule ...Rule) *FutureEvent {
	return &FutureEvent{
		Priority: Priority,
		Rule:     rule,
		Block:    Block,
	}
}

func (m *Matcher) FutureEvent(rule ...Rule) *FutureEvent {
	return &FutureEvent{
		Priority: m.Priority,
		Block:    m.Block,
		Rule:     rule,
	}
}

// Next 返回一个 chan 用于接收下一个指定事件
// 该 chan 必须接收，如需手动取消监听，请使用 Repeat 方法
func (n *FutureEvent) Next() <-chan *Ctx {
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
func (n *FutureEvent) Repeat() (recv <-chan *Ctx, cancel func()) {
	ch, done := make(chan *Ctx, 1), make(chan struct{})
	go func() {
		defer close(ch)
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
				return
			}
		}
	}()
	return ch, func() {
		close(done)
	}
}

// Take 基于 Repeat 封装，返回一个 chan 接收指定数量的事件
// 该 chan 对象必须接收，否则将有 goroutine 泄漏，如需手动取消请使用 Repeat
func (n *FutureEvent) Take(num int) <-chan *Ctx {
	recv, cancel := n.Repeat()
	ch := make(chan *Ctx, num)
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			ch <- <-recv
		}
		cancel()
	}()
	return ch
}
