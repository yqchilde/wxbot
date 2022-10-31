package weather

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type Weather struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Key    string `yaml:"key"`
}

var (
	pluginInfo = &Weather{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {XX天气} => 获取天气数据，Ps:济南天气、北京-朝阳天气",
			Commands: []string{`([^\x00-\xff]{0,6}-?[^\x00-\xff]{0,6})天气`},
			Weight:   98,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (m *Weather) OnRegister() {}

func (m *Weather) OnEvent(msg *robot.Message) {
	if msg != nil {
		if idx, ok := msg.MatchRegexCommand(pluginInfo.Commands); ok {
			var re = regexp.MustCompile(pluginInfo.Commands[idx])
			match := re.FindAllStringSubmatch(msg.Content.Msg, -1)
			city := match[0][1]
			apiKey := plugin.RawConfig.Get("key").(string)
			locationSplit := strings.Split(city, "-")
			var locationList []Location
			if len(locationSplit) == 1 {
				locationList = getCityLocation(apiKey, "", locationSplit[0])
			}
			if len(locationSplit) == 2 {
				locationList = getCityLocation(apiKey, locationSplit[0], locationSplit[1])
			}
			if len(locationList) == 0 {
				msg.ReplyTextAndAt("未找到城市")
				return
			} else if len(locationList) == 1 {
				location = locationList[0].Id
			} else {
				adm := map[string]struct{}{}
				for _, v := range locationList {
					adm[v.Adm2] = struct{}{}
				}
				if len(adm) == 1 {
					location = locationList[0].Id
				} else {
					msg.ReplyTextAndAt("查询到多个地区地址，请输入更详细的地区，比如：北京-朝阳天气")
					return
				}
			}

			weatherNow := getWeatherNow(apiKey, location)
			weather2d := getWeather2d(apiKey, location)
			weatherIndices := getWeatherIndices(apiKey, location)
			console := "城市: %s\n"
			console += "今天: %s\n"
			console += "当前温度: %s°，体感温度: %s°\n"
			console += "白天: %s(%s°-%s°)，夜间: %s\n"
			console += "日出时间: %s，日落时间: %s\n"
			console += "当前降水量: %s，能见度: %s，云量: %s\n"
			console += "天气舒适指数: %s\n"
			console += "\n"
			console += "明天: %s\n"
			console += "白天: %s(%s°-%s°)，夜间: %s\n"
			console += "日出时间: %s，日落时间: %s\n"
			console = fmt.Sprintf(console, locationList[0].Name, weather2d[0].FxDate, weatherNow.Temp, weatherNow.FeelsLike, weather2d[0].TextDay, weather2d[0].TempMin, weather2d[0].TempMax, weather2d[0].TextNight, weather2d[0].Sunrise, weather2d[0].Sunset, weatherNow.Precip, weatherNow.Vis, weatherNow.Cloud, weatherIndices, weather2d[1].FxDate, weather2d[1].TextDay, weather2d[1].TempMin, weather2d[1].TempMax, weather2d[1].TextNight, weather2d[1].Sunrise, weather2d[1].Sunset)
			msg.ReplyText(console)
		}
	}
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

// 城市搜索
func getCityLocation(key, adm, location string) []Location {
	resp := req.C().Get("https://geoapi.qweather.com/v2/city/lookup").
		SetQueryParams(map[string]string{
			"key":      key,
			"adm":      adm,
			"location": location,
		}).Do()

	var locationList []Location
	gjson.Get(resp.String(), "location").ForEach(func(key, value gjson.Result) bool {
		locationList = append(locationList, Location{
			Name: value.Get("name").String(),
			Id:   value.Get("id").String(),
			Adm2: value.Get("adm2").String(),
			Adm1: value.Get("adm1").String(),
		})
		return true
	})
	return locationList
}

// 实时天气
func getWeatherNow(key, location string) WeatherNow {
	resp := req.C().Get("https://devapi.qweather.com/v7/weather/now").
		SetQueryParams(map[string]string{
			"key":      key,
			"location": location,
		}).Do()

	data := gjson.Get(resp.String(), "now")
	return WeatherNow{
		UpdateTime: data.Get("obsTime").String(),
		Temp:       data.Get("temp").String(),
		FeelsLike:  data.Get("feelsLike").String(),
		Text:       data.Get("text").String(),
		Precip:     data.Get("precip").String(),
		Vis:        data.Get("vis").String(),
		Cloud:      data.Get("cloud").String(),
	}
}

// 两天天气预报
func getWeather2d(key, location string) []WeatherDay {
	resp := req.C().Get("https://devapi.qweather.com/v7/weather/3d").
		SetQueryParams(map[string]string{
			"key":      key,
			"location": location,
		}).Do()

	data := gjson.Get(resp.String(), "daily")
	return []WeatherDay{
		{
			FxDate:    data.Get("0.fxDate").String(),
			Sunrise:   data.Get("0.sunrise").String(),
			Sunset:    data.Get("0.sunset").String(),
			TempMax:   data.Get("0.tempMax").String(),
			TempMin:   data.Get("0.tempMin").String(),
			TextDay:   data.Get("0.textDay").String(),
			TextNight: data.Get("0.textNight").String(),
		},
		{
			FxDate:    data.Get("1.fxDate").String(),
			Sunrise:   data.Get("1.sunrise").String(),
			Sunset:    data.Get("1.sunset").String(),
			TempMax:   data.Get("1.tempMax").String(),
			TempMin:   data.Get("1.tempMin").String(),
			TextDay:   data.Get("1.textDay").String(),
			TextNight: data.Get("1.textNight").String(),
		},
	}
}

// 分钟级降水
func getMinutely5m(key, location string) {
	resp := req.C().Get("https://devapi.qweather.com/v7/minutely/5m").
		SetQueryParams(map[string]string{
			"key":      key,
			"location": location,
		}).Do()

	_ = gjson.Get(resp.String(), "minutely")
}

// 天气指数
func getWeatherIndices(key, location string) string {
	resp := req.C().Get("https://devapi.qweather.com/v7/indices/1d").
		SetQueryParams(map[string]string{
			"key":      key,
			"location": location,
			"type":     "8",
		}).Do()

	return gjson.Get(resp.String(), "daily").Get("0.text").String()
}
