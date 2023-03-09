package manager

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/yqchilde/wxbot/engine/robot"
)

var job = gocron.NewScheduler(time.Local)

func init() {
	// 设置最大并发任务数，防止并发太大对账号有影响
	// SetMaxConcurrentJobs函数可能有重复的可能，待观察，https://github.com/go-co-op/gocron/blob/main/executor.go#L16
	job.SetMaxConcurrentJobs(10, gocron.WaitMode)
}

// AddCronjobOfEveryMonth 添加每月提醒
func AddCronjobOfEveryMonth(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	timeSplit := strings.Split(matched[2], ":")
	hour, minute, second := timeSplit[0], timeSplit[1], timeSplit[2]
	taskCron := fmt.Sprintf("%s %s %s %s * *", second, minute, hour, matched[1])
	return job.CronWithSeconds(taskCron).Tag(jobTag).Do(func() { f() })
}

// AddCronjobOfEveryWeek 添加每周提醒
func AddCronjobOfEveryWeek(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
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

// AddCronjobOfEveryDay 添加每天提醒
func AddCronjobOfEveryDay(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	return job.Every(1).Day().At(matched[1]).Tag(jobTag).Do(func() { f() })
}

// AddCronjobForInterval 添加间隔提醒
func AddCronjobForInterval(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
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

// AddCronjobForSpecifyTime 添加指定时间提醒
func AddCronjobForSpecifyTime(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	parseTime, _ := time.ParseInLocation("2006-01-02 15:04:05", matched[1], time.Local)
	if parseTime.Before(time.Now()) {
		return nil, fmt.Errorf("请不要设置过去的时间")
	}
	return job.Every(1).LimitRunsTo(1).StartAt(parseTime).Tag(jobTag).Do(func() {
		f()
		db.Orm.Table("cronjob").Where("tag = ?", jobTag).Delete(&CronJob{})
	})
}

// AddCronjobForExpression 添加表达式提醒
func AddCronjobForExpression(ctx *robot.Ctx, jobTag string, matched []string, f func()) (*gocron.Job, error) {
	return job.CronWithSeconds(matched[1]).Tag(jobTag).Do(func() { f() })
}
