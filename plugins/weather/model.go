package weather

type Weather struct {
	AppKey string `gorm:"column:app_key"`
}

var location string

type Location struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Adm2 string `json:"adm2"`
	Adm1 string `json:"adm1"`
}

type WeatherNow struct {
	UpdateTime string `json:"updateTime"` // 更新时间
	Temp       string `json:"temp"`       // 温度
	FeelsLike  string `json:"feelsLike"`  // 体感温度
	Text       string `json:"text"`       // 天气状况
	Precip     string `json:"precip"`     // 降水量
	Vis        string `json:"vis"`        // 能见度
	Cloud      string `json:"cloud"`      // 云量
}

type WeatherDay struct {
	FxDate    string `json:"fxDate"`    // 预报日期
	Sunrise   string `json:"sunrise"`   // 日出时间
	Sunset    string `json:"sunset"`    // 日落时间
	TempMax   string `json:"tempMax"`   // 最高温度
	TempMin   string `json:"tempMin"`   // 最低温度
	TextDay   string `json:"textDay"`   // 白天天气现象文字
	TextNight string `json:"textNight"` // 晚间天气现象文字
}
