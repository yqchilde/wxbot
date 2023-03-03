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

// fileSecret 文件服务秘钥
var fileSecret []byte

// SetFileSecret 设置文件服务秘钥
func SetFileSecret(secret []byte) {
	fileSecret = secret
}

// GetFileSecret 获取文件服务秘钥
func (ctx *Ctx) GetFileSecret() []byte {
	return fileSecret
}
