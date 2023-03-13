package control

import "github.com/yqchilde/wxbot/engine/robot"

type Matcher robot.Matcher

// SetBlock 设置是否阻断后面的Matcher触发
func (m *Matcher) SetBlock(block bool) *Matcher {
	_ = (*robot.Matcher)(m).SetBlock(block)
	return m
}

// SetPriority 设置当前Matcher优先级
func (m *Matcher) SetPriority(priority uint64) *Matcher {
	_ = (*robot.Matcher)(m).SetPriority(priority)
	return m
}

// Handle 直接处理事件
func (m *Matcher) Handle(handler robot.Handler) {
	_ = (*robot.Matcher)(m).Handle(handler)
}
