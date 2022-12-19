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

// GetOptions 获取当前群组/私聊插件控制器配置
func GetOptions(wxId string) *MenuOptions {
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
