package control

import (
	"sync/atomic"

	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	priority uint64
)

// Register 注册插件控制器
func Register(service string, o *Options[*robot.Ctx]) *Engine {
	engine := newEngine(service, int(atomic.AddUint64(&priority, 10)), o)
	return engine
}
