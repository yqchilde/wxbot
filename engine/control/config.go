package control

import "github.com/yqchilde/wxbot/engine/robot"

type Options struct {
	Alias            string               // 插件别名
	Help             string               // 插件帮助信息
	Priority         uint64               // 优先级,只读
	DisableOnDefault bool                 // 默认禁用状态
	HideMenu         bool                 // 是否隐藏在菜单中，默认显示，可用于隐藏一些不希望展示的插件
	DataFolder       string               // 数据文件夹
	OnEnable         func(ctx *robot.Ctx) // 自定义启用插件后执行的操作
	OnDisable        func(ctx *robot.Ctx) // 自定义禁用插件后执行的操作
}

// PluginConfig 插件配置表
type PluginConfig struct {
	GroupID string `gorm:"column:gid"`    // 群组ID
	Enable  bool   `gorm:"column:enable"` // 启用状态，默认启用
}

// PluginBanConfig 插件Ban配置表
type PluginBanConfig struct {
	Label   string `gorm:"column:label"` // 标签(service_uid_gid)
	UserID  string `gorm:"column:uid"`   // 用户ID
	GroupID string `gorm:"column:gid"`   // 群组ID
}

// BotBlockConfig bot ban配置表，ban掉的用户无法使用所有插件
type BotBlockConfig struct {
	UserID string `gorm:"column:uid"` // 用户ID
}

// BotResponseConfig bot响应群配置表
type BotResponseConfig struct {
	GroupID string `gorm:"column:gid"`    // 群组ID
	Status  bool   `gorm:"column:status"` // 响应状态，默认启用
}

// MessageRecord 消息记录表
type MessageRecord struct {
	*robot.MessageRecord
}

// BaseConfig 基础配置表
type BaseConfig struct {
	FileSecret []byte `gorm:"column:file_secret"`
}
