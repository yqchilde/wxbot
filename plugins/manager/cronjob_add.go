package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/yqchilde/wxbot/engine/robot"
)

var job = gocron.NewScheduler(time.Local)

// AddRemindOfEveryMonth 添加每月提醒
func AddRemindOfEveryMonth(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	timeSplit := strings.Split(matched[2], ":")
	hour, minute, second := timeSplit[0], timeSplit[1], timeSplit[2]
	taskCron := fmt.Sprintf("%s %s %s %s * *", second, minute, hour, matched[1])
	return job.Tag(jobTag).CronWithSeconds(taskCron).Do(func() { f() })
}

// AddRemindOfEveryWeek 添加每周提醒
func AddRemindOfEveryWeek(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	week, timed := matched[1], matched[2]
	switch week {
	case "一":
		job.Every(1).Monday().At(timed)
	case "二":
		job.Every(1).Tuesday().At(timed)
	case "三":
		job.Every(1).Wednesday().At(timed)
	case "四":
		job.Every(1).Thursday().At(timed)
	case "五":
		job.Every(1).Friday().At(timed)
	case "六":
		job.Every(1).Saturday().At(timed)
	case "七", "日":
		job.Every(1).Sunday().At(timed)
	}
	return job.Tag(jobTag).Do(func() { f() })
}

// AddRemindOfEveryDay 添加每天提醒
func AddRemindOfEveryDay(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	return job.Tag(jobTag).Every(1).Day().At(matched[1]).Do(func() { f() })
}

// AddRemindForInterval 添加间隔提醒
func AddRemindForInterval(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	duration, unit := matched[1], matched[2]
	switch unit {
	case "秒", "s":
		if dur, err := time.ParseDuration(fmt.Sprintf("%ss", duration)); err == nil {
			job.Every(dur).StartAt(time.Now().Add(dur))
		}
	case "分", "分钟", "m":
		if dur, err := time.ParseDuration(fmt.Sprintf("%sm", duration)); err == nil {
			job.Every(dur).StartAt(time.Now().Add(dur))
		}
	case "时", "小时", "h":
		if dur, err := time.ParseDuration(fmt.Sprintf("%sh", duration)); err == nil {
			job.Every(dur).StartAt(time.Now().Add(dur))
		}
	}
	return job.Tag(jobTag).Do(func() { f() })
}
