package weather

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db      sqlite.DB
	weather Weather
)

func init() {
	engine := control.Register("weather", &control.Options{
		Alias: "天气查询",
		Help: "指令:\n" +
			"* [城市名]天气 -> 获取天气数据，Ps:济南天气、北京-朝阳天气",
		DataFolder: "weather",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/weather.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("weather", &weather); err != nil {
		log.Fatalf("create weather table failed: %v", err)
	}

	engine.OnRegex(`^([^\x00-\xff]{2,6}-?[^\x00-\xff]{0,6})天气$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if weather.AppKey == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置appKey\n指令：set weather appKey __\n相关秘钥申请地址：https://dev.qweather.com")
			return
		}

		city := ctx.State["regex_matched"].([]string)[1]
		locationSplit := strings.Split(city, "-")
		var locationList []Location
		if len(locationSplit) == 1 {
			locationList = getCityLocation(weather.AppKey, "", locationSplit[0])
		}
		if len(locationSplit) == 2 {
			locationList = getCityLocation(weather.AppKey, locationSplit[0], locationSplit[1])
		}
		if len(locationList) == 0 {
			ctx.ReplyTextAndAt("未找到城市")
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
				ctx.ReplyTextAndAt("查询到多个地区地址，请输入更详细的地区，比如：北京-朝阳天气")
				return
			}
		}

		weatherNow := getWeatherNow(weather.AppKey, location)
		weather2d := getWeather2d(weather.AppKey, location)
		weatherIndices := getWeatherIndices(weather.AppKey, location)
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
		console = fmt.Sprintf(console, locationList[0].Adm1+"/"+locationList[0].Adm2+"/"+locationList[0].Name, weather2d[0].FxDate, weatherNow.Temp, weatherNow.FeelsLike, weather2d[0].TextDay, weather2d[0].TempMin, weather2d[0].TempMax, weather2d[0].TextNight, weather2d[0].Sunrise, weather2d[0].Sunset, weatherNow.Precip, weatherNow.Vis, weatherNow.Cloud, weatherIndices, weather2d[1].FxDate, weather2d[1].TextDay, weather2d[1].TempMin, weather2d[1].TempMax, weather2d[1].TextNight, weather2d[1].Sunrise, weather2d[1].Sunset)
		ctx.ReplyText(console)
	})

	engine.OnRegex("set weather appKey ([0-9a-z]{32})", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		appKey := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("weather").Where("1 = 1").Update("app_key", appKey).Error; err != nil {
			ctx.ReplyTextAndAt("appKey配置失败")
			return
		}
		weather.AppKey = appKey
		ctx.ReplyText("appKey设置成功")
	})

	engine.OnFullMatch("get weather info", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var weather Weather
		if err := db.Orm.Table("weather").Limit(1).Find(&weather).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - 查询天气\nappKey: %s", weather.AppKey))
	})
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
