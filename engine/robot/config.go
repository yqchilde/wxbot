package robot

import "time"

type Config struct {
	BotWxId        string        `mapstructure:"botWxId"`       // 机器人微信ID
	BotNickname    string        `mapstructure:"botNickname"`   // 机器人名称
	SuperUsers     []string      `mapstructure:"superUsers"`    // 超级用户(管理员)
	CommandPrefix  string        `mapstructure:"commandPrefix"` // 管理员触发命令
	ServerPort     uint          `mapstructure:"serverPort"`    // 启动HTTP服务端口
	BufferLen      uint          `mapstructure:"-"`             // 事件缓冲区长度, 默认4096
	Latency        time.Duration `mapstructure:"-"`             // 事件处理延迟 (延迟 latency + (0~100ms) 再处理事件) (默认1s)
	MaxProcessTime time.Duration `mapstructure:"-"`             // 事件最大处理时间 (默认3min)
	Framework      struct {
		Name     string `mapstructure:"name"`     // 接入框架名称
		ApiUrl   string `mapstructure:"apiUrl"`   // 接入框架API地址
		ApiToken string `mapstructure:"apiToken"` // 接入框架API Token
	} `mapstructure:"framework"`
}
