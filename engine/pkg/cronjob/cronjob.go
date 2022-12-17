package cronjob

func ParseToCron(duration string, unit string) string {
	switch unit {
	case "秒", "s":
		return "@every " + duration + "s"
	case "分", "分钟", "m":
		return "@every " + duration + "m"
	case "时", "小时", "h":
		return "@every " + duration + "h"
	case "天", "日", "d":
		return "@every " + duration + "d"
	}
	return ""
}
