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

// GetOptionsMenu 获取当前群组/私聊插件控制器菜单配置
func GetOptionsMenu(wxId string) *MenuOptions {
	services := managers.LookupAll()
	menuOptions := MenuOptions{WxId: wxId}
	for _, s := range services {
		if s.Options.HideMenu {
			continue
		}
		menuOptions.Menus = append(menuOptions.Menus, struct {
			Name      string `json:"name"`
			Alias     string `json:"alias"`
			Priority  int    `json:"priority"`
			Describe  string `json:"describe"`
			DefStatus bool   `json:"defStatus"`
			CurStatus bool   `json:"curStatus"`
		}{
			Name:      s.Service,
			Alias:     s.Options.Alias,
			Priority:  s.Options.priority,
			Describe:  s.Options.Help,
			DefStatus: !s.Options.DisableOnDefault,
			CurStatus: s.IsEnabledIn(wxId),
		})
	}
	return &menuOptions
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
