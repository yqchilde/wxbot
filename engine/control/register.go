package control

import (
	"sync/atomic"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	priority    uint64
	priorityMap = make(map[uint64]string)
)

// Register 注册插件控制器
func Register(service string, o *Options[*robot.Ctx]) *Engine {
	atomic.AddUint64(&priority, 10)
	s, ok := priorityMap[priority]
	if ok {
		log.Fatalf("[%s]插件优先级 %d 已被 %s 占用", service, priority, s)
	}
	priorityMap[priority] = service
	log.Printf("[%s]插件已注册, 优先级: %d", service, priority)
	return newEngine(service, o)
}

// GetOptionsOnCronjob 获取定时任务插件控制器配置
func GetOptionsOnCronjob() map[string]*Control[*robot.Ctx] {
	var (
		services      = managers.LookupAll()
		servicesClone = make(map[string]*Control[*robot.Ctx])
	)

	for i := range services {
		if services[i].Options.OnCronjob == nil {
			continue
		}
		servicesClone[services[i].Service] = services[i]
	}
	return servicesClone
}
