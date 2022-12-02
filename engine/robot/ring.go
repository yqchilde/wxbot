package robot

import (
	"container/ring"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type eventRing struct {
	sync.Mutex
	r *ring.Ring
}

type eventRingItem struct {
	response []byte
	caller   APICaller
}

func newRing(ringLen uint) eventRing {
	n := int(ringLen)
	r := ring.New(n)
	// Initialize the ring with locked eventRing
	for i := 0; i < n; i++ {
		r.Value = (*eventRingItem)(nil)
		r = r.Next()
	}
	return eventRing{r: r}
}

// processEvent 同步向池中放入事件
func (evr *eventRing) processEvent(response []byte, caller APICaller) {
	evr.Lock()
	defer evr.Unlock()
	r := evr.r
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&r.Value), unsafe.Sizeof(uintptr(0)))),
		unsafe.Pointer(&eventRingItem{
			response: response,
			caller:   caller,
		}),
	)
	evr.r = r.Next()
}

// loop 循环处理事件
//
//	latency 延迟 latency 再处理事件
func (evr *eventRing) loop(latency, maxWait time.Duration, process func([]byte, APICaller, time.Duration)) {
	go func(r *ring.Ring) {
		for range time.NewTicker(latency).C {
			it := r.Value.(*eventRingItem)
			if it == nil { // 还未有消息
				continue
			}
			process(it.response, it.caller, maxWait)
			it.response = nil
			it.caller = nil
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&r.Value), unsafe.Sizeof(uintptr(0)))), unsafe.Pointer(nil))
			r = r.Next()
			runtime.GC()
		}
	}(evr.r)
}
