package robot

import (
	"sync"
	"time"
)

type EventBuffer struct {
	sync.Mutex
	items chan EventBufferItem
	done  chan struct{}
}

type EventBufferItem struct {
	event     *Event
	framework IFramework
}

func NewEventBuffer(bufferLen uint) *EventBuffer {
	return &EventBuffer{
		items: make(chan EventBufferItem, bufferLen),
		done:  make(chan struct{}),
	}
}

// ProcessEvent 处理事件
func (evr *EventBuffer) ProcessEvent(event *Event, framework IFramework) {
	evr.items <- EventBufferItem{
		event:     event,
		framework: framework,
	}
}

// Loop 以给定的延迟和最长等待时间处理环中的事件
func (evr *EventBuffer) Loop(latency, maxWait time.Duration, process func(*Event, IFramework, time.Duration)) {
	go func() {
		ticker := time.NewTicker(latency)
		for {
			select {
			case item := <-evr.items:
				process(item.event, item.framework, maxWait)
			case <-ticker.C:
			case <-evr.done:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止事件处理循环
func (evr *EventBuffer) Stop() {
	evr.done <- struct{}{}
}
