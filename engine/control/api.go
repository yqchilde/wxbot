package control

import "github.com/yqchilde/wxbot/engine/robot"

type controlApi struct{}

// GetMenus 获取菜单
func (c *controlApi) GetMenus(wxId string) (menus []map[string]interface{}) {
	// 检查wxId是否存在
	users := robot.GetBot().Users()
	for _, v := range users {
		if v.WxId == wxId {
			for _, v := range managers.LookupAll() {
				if v.Options.HideMenu {
					continue
				}
				menus = append(menus, map[string]interface{}{
					"name":      v.Service,
					"alias":     v.Options.Alias,
					"priority":  v.Options.Priority,
					"describe":  v.Options.Help,
					"defStatus": !v.Options.DisableOnDefault,
					"curStatus": v.IsEnabledIn(wxId),
				})
			}
			break
		}
	}
	return
}
