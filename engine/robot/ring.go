package robot

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type eventRing struct {
	sync.Mutex
	c uintptr
	r []*eventRingItem
	i uintptr
	p []eventRingItem
}

type eventRingItem struct {
	event     *Event
	framework IFramework
}

func newRing(ringLen uint) eventRing {
	return eventRing{
		r: make([]*eventRingItem, ringLen),
		p: make([]eventRingItem, ringLen+1),
	}
}

// processEvent 同步向池中放入事件
func (evr *eventRing) processEvent(event *Event, framework IFramework) {
	evr.Lock()
	defer evr.Unlock()
	r := evr.c % uintptr(len(evr.r))
	p := evr.i % uintptr(len(evr.p))
	evr.p[p] = eventRingItem{
		event:     event,
		framework: framework,
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&evr.r[r])), unsafe.Pointer(&evr.p[p]))
	evr.c++
	evr.i++
}

// loop 循环处理事件，latency 延迟 latency 再处理事件
func (evr *eventRing) loop(latency, maxWait time.Duration, process func(*Event, IFramework, time.Duration)) {
	go func(r []*eventRingItem) {
		c := uintptr(0)
		for range time.NewTicker(latency).C {
			i := c % uintptr(len(r))
			it := (*eventRingItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r[i]))))
			if it == nil { // 还未有消息
				continue
			}
			process(it.event, it.framework, maxWait)
			it.event = nil
			it.framework = nil
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r[i])), unsafe.Pointer(nil))
			c++
			runtime.GC()
		}
	}(evr.r)
}
