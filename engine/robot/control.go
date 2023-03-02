package robot

// ControlApi 可获取control的api
var ControlApi controlApi

type controlApi interface {
	// GetMenus 获取菜单
	GetMenus(wxId string) []map[string]interface{}
}

// RegisterApi 注册api
func RegisterApi(api controlApi) {
	ControlApi = api
}
