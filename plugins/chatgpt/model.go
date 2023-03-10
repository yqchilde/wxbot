package chatgpt

// ApiKey 表名:apikey，存放openai key
type ApiKey struct {
	Key string `gorm:"column:key;index"`
}

// ApiProxy 表名:apiproxy，存放openai 代理url地址
type ApiProxy struct {
	Id  uint   `gorm:"column:id;index"`
	Url string `gorm:"column:url;"`
}

// GptModel 表名:gptmodel，存放gpt模型相关配置参数
type GptModel struct {
	Model            string  `gorm:"column:model"`
	MaxTokens        int     `gorm:"column:max_tokens"`
	Temperature      float64 `gorm:"column:temperature"`
	TopP             float64 `gorm:"column:top_p"`
	PresencePenalty  float64 `gorm:"column:presence_penalty"`
	FrequencyPenalty float64 `gorm:"column:frequency_penalty"`
	ImageSize        string  `gorm:"column:image_size"`
}

var defaultGptModel = GptModel{
	Model:            "gpt-3.5-turbo",
	MaxTokens:        4096,
	Temperature:      0.8,
	TopP:             1.0,
	PresencePenalty:  0.0,
	FrequencyPenalty: 0.6,
	ImageSize:        "512x512",
}

// SystemRoles 表名:roles，存放系统角色
type SystemRoles struct {
	Role string `gorm:"column:role"`
	Desc string `gorm:"column:desc"`
}

// SensitiveWords 表名:sensitive，存放敏感词
type SensitiveWords struct {
	Type    int    `gorm:"column:type;index"`    // 1:内置敏感词，2:自定义敏感词
	Word    string `gorm:"column:word;index"`    // 敏感词
	Deleted int    `gorm:"column:deleted;index"` // 0:未删除，1:已删除
}
