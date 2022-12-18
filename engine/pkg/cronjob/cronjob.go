package cronjob

import (
	"fmt"
	"strconv"
)

func ParseToCron(duration string, unit string) string {
	switch unit {
	case "秒", "s":
		return fmt.Sprintf("@every %ss", duration)
	case "分", "分钟", "m":
		return fmt.Sprintf("@every %sm", duration)
	case "时", "小时", "h":
		return fmt.Sprintf("@every %sh", duration)
	case "天", "日", "d":
		if durInt, err := strconv.Atoi(duration); err != nil {
			return ""
		} else {
			return fmt.Sprintf("@every %dh", durInt*24)
		}
	}
	return ""
}
