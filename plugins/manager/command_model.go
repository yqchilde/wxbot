package manager

// MenuOptions 菜单配置
type MenuOptions struct {
	WxId  string `json:"wxId"`
	Menus []struct {
		Name      string `json:"name"`
		Alias     string `json:"alias"`
		Priority  uint64 `json:"priority"`
		Describe  string `json:"describe"`
		DefStatus bool   `json:"defStatus"`
		CurStatus bool   `json:"curStatus"`
	} `json:"menus"`
}
