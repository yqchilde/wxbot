package manager

import (
	"embed"
	"fmt"
	"regexp"
	"time"

	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

//go:embed data
var holidayData embed.FS

const (
	JobTypeRemind = "remind" // 提醒类任务
	JobTypePlugin = "plugin" // 插件类任务
)

type CronJob struct {
	Id     int64  `gorm:"primary_key"`    // 任务ID
	Tag    string `gorm:"column:tag"`     // 任务标签
	Type   string `gorm:"column:type"`    // 任务类型
	Desc   string `gorm:"column:desc"`    // 任务描述
	WxId   string `gorm:"column:wx_id"`   // 微信ID
	WxName string `gorm:"column:wx_name"` // 微信昵称
	Remind string `gorm:"column:remind"`  // 提醒内容
}

func registerCronjob() {
	engine := control.Register("cronjob", &control.Options{
		Alias: "定时任务",
		Help: "权限:\n" +
			"仅限机器人管理员\n\n" +
			"指令:\n" +
			"提醒类任务指令:\n" +
			"* 设置每月[]号[]的提醒任务 -> 例如：设置每月8号10:00:00的提醒任务\n" +
			"* 设置每周[][]的提醒任务 -> 例如：设置每周三10:00:00的提醒任务\n" +
			"* 设置每天[]的提醒任务 -> 例如：设置每天10:00:00的提醒任务\n" +
			"* 设置每隔[]的提醒任务 -> 例如：设置每隔1小时的提醒任务\n" +
			"* 设置[]的提醒任务 -> 例如：设置2023-01-01 15:00:00的提醒任务\n" +
			"* 设置表达式[]的提醒任务 -> 例如：设置表达式(*/10 * * * * *)的提醒任务\n" +
			"* 设置工作日[]的提醒任务 -> 例如：设置工作日10:00:00的提醒任务\n" +
			"* 删除全部提醒任务\n\n" +
			"插件类任务指令:\n" +
			"* 设置每月[]号[]的插件任务 -> 例如：设置每月8号10:00:00的插件任务\n" +
			"* 设置每周[][]的插件任务 -> 例如：设置每周三10:00:00的插件任务\n" +
			"* 设置每天[]的插件任务 -> 例如：设置每天10:00:00的插件任务\n" +
			"* 设置每隔[]的插件任务 -> 例如：设置每隔1小时的插件任务\n" +
			"* 设置[]的插件任务 -> 例如：设置2023-01-01 15:00:00的插件任务\n" +
			"* 设置表达式[]的插件任务 -> 例如：设置表达式(*/10 * * * * *)的插件任务\n" +
			"* 设置工作日[]的插件任务 -> 例如：设置工作日10:00:00的插件任务\n" +
			"* 删除全部插件任务\n\n" +
			"其他指令:\n" +
			"* 列出所有任务\n" +
			"* 删除任务 [任务ID]\n" +
			"* 删除全部任务\n",
		DataFolder: "manager",
	})
	if err := db.Create("cronjob", &CronJob{}); err != nil {
		log.Fatalf("create cronjob table failed: %v", err)
	}

	go func() {
		// 用于恢复定时任务
		var cronJobs []CronJob
		if err := db.Orm.Table("cronjob").Find(&cronJobs).Error; err != nil {
			return
		}
		ctx := robot.GetCtx()
		for i := range cronJobs {
			cronJob := cronJobs[i]
			switch cronJob.Type {
			case JobTypeRemind:
				// 恢复每月的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryMonth).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryMonth(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每月提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每周的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryWeek).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryWeek(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每周提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每天的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryDay(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每天提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复间隔提醒任务
				if matched := regexp.MustCompile(RegexOfRemindInterval).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForInterval(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复间隔提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复指定时间提醒任务
				if matched := regexp.MustCompile(RegexOfRemindSpecifyTime).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForSpecifyTime(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复指定时间提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复表达式提醒任务
				if matched := regexp.MustCompile(RegexOfRemindExpression).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForExpression(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复表达式提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复工作日提醒任务
				if matched := regexp.MustCompile(RegexOfRemindWorkDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryDay(ctx, cronJob.Tag, matched, func() {
						data, err := holidayData.ReadFile(fmt.Sprintf("data/holiday_%d.json", time.Now().Year()))
						if err != nil {
							log.Errorf("获取节假日数据失败: %v", err)
							ctx.ReplyText("很抱歉今天提醒您，由于节假日数据获取失败了，我无法确定今天是否是工作日")
							return
						}
						var now = time.Now().Local()
						var isWorkDay, isHoliday bool
						gjson.GetBytes(data, "days").ForEach(func(key, val gjson.Result) bool {
							if val.Get("date").String() == now.Format("2006-01-02") {
								isWorkDay = !val.Get("isOffDay").Bool()
								isHoliday = val.Get("isOffDay").Bool()
								return false
							}
							return true
						})
						if !isHoliday && !isWorkDay {
							if now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
								isWorkDay = true
							}
						}
						if isWorkDay {
							ctx.SendText(cronJob.WxId, cronJob.Remind)
						}
					}); err != nil {
						log.Errorf("恢复表达式提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}
			case JobTypePlugin:
				// 恢复每月的插件任务
				if matched := regexp.MustCompile(RegexOfPluginEveryMonth).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryMonth(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每月插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每周的插件任务
				if matched := regexp.MustCompile(RegexOfPluginEveryWeek).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryWeek(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每周插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每天的插件任务
				if matched := regexp.MustCompile(RegexOfPluginEveryDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryDay(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每天插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复间隔插件任务
				if matched := regexp.MustCompile(RegexOfPluginInterval).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForInterval(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复间隔插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复指定时间插件任务
				if matched := regexp.MustCompile(RegexOfPluginSpecifyTime).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForSpecifyTime(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复指定时间插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复表达式插件任务
				if matched := regexp.MustCompile(RegexOfPluginExpression).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobForExpression(ctx, cronJob.Tag, matched, func() {
						ctx.SendEvent(cronJob.WxId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复表达式插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复工作日插件任务
				if matched := regexp.MustCompile(RegexOfPluginWorkDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddCronjobOfEveryDay(ctx, cronJob.Tag, matched, func() {
						data, err := holidayData.ReadFile(fmt.Sprintf("data/holiday_%d.json", time.Now().Year()))
						if err != nil {
							log.Errorf("获取节假日数据失败: %v", err)
							ctx.SendText(cronJob.WxId, "很抱歉今天提醒您，由于节假日数据获取失败了，我无法确定今天是否是工作日")
							return
						}
						var now = time.Now().Local()
						var isWorkDay, isHoliday bool
						gjson.GetBytes(data, "days").ForEach(func(key, val gjson.Result) bool {
							if val.Get("date").String() == now.Format("2006-01-02") {
								isWorkDay = !val.Get("isOffDay").Bool()
								isHoliday = val.Get("isOffDay").Bool()
								return false
							}
							return true
						})
						if !isHoliday && !isWorkDay {
							if now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
								isWorkDay = true
							}
						}
						if isWorkDay {
							ctx.SendEvent(cronJob.WxId, cronJob.Remind)
						}
					}); err != nil {
						log.Errorf("恢复表达式提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}
			}
		}
		job.StartAsync()
	}()

	// 设置提醒类任务指令
	SetRemindCommand(engine)

	// 设置插件类任务指令
	SetPluginCommand(engine)

	// 列出当前所有定时任务
	engine.OnFullMatch("列出所有任务").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var cronJobs []CronJob
		if err := db.Orm.Table("cronjob").Where("wx_id = ?", ctx.Event.FromUniqueID).Find(&cronJobs).Error; err != nil {
			ctx.ReplyTextAndAt("查询定时任务失败")
			return
		}
		var jobInfo string
		for i := range cronJobs {
			switch cronJobs[i].Type {
			case JobTypeRemind:
				jobInfo += fmt.Sprintf("任务ID: %d\n任务类型: %s\n任务描述: %s\n任务内容: %s\n\n", cronJobs[i].Id, cronJobs[i].Type, cronJobs[i].Desc, cronJobs[i].Remind)
			case JobTypePlugin:
				jobInfo += fmt.Sprintf("任务ID: %d\n任务类型: %s\n任务描述: %s\n任务内容: %s\n\n", cronJobs[i].Id, cronJobs[i].Type, cronJobs[i].Desc, cronJobs[i].Remind)
			}
		}
		if len(cronJobs) == 0 {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务", len(cronJobs)))
		} else {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务:\n%s", len(cronJobs), jobInfo))
		}
	})

	// 删除任务 任务ID
	engine.OnRegex(`^删除任务 ?(\d+)$`, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		jobId := ctx.State["regex_matched"].([]string)[1]
		var jobTag string
		if err := db.Orm.Table("cronjob").Where("id = ?", jobId).Pluck("tag", &jobTag).Error; err != nil {
			log.Errorf("[CronJob] 删除任务失败: %v", err)
			ctx.ReplyTextAndAt("删除任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("id = ?", jobId).Delete(&CronJob{}).Error; err != nil {
			log.Errorf("[CronJob] 删除任务失败: %v", err)
			ctx.ReplyTextAndAt("删除任务失败")
			return
		}
		if err := job.RemoveByTag(jobTag); err != nil {
			log.Errorf("[CronJob] 删除任务失败: %v", err)
		} else {
			ctx.ReplyTextAndAt(fmt.Sprintf("任务[%s]删除成功", jobId))
		}
	})

	// 删除全部任务
	engine.OnFullMatchGroup([]string{"删除全部任务", "删除所有任务"}, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jobTags []string
		if err := db.Orm.Table("cronjob").Where("wx_id = ?", ctx.Event.FromUniqueID).Pluck("tag", &jobTags).Error; err != nil {
			log.Errorf("[CronJob] 删除全部任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("wx_id = ?", ctx.Event.FromUniqueID).Delete(&CronJob{}).Error; err != nil {
			log.Errorf("[CronJob] 删除全部任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部任务失败")
			return
		}
		if err := job.RemoveByTagsAny(jobTags...); err != nil {
			log.Errorf("[CronJob] 删除全部任务失败: %v", err)
		} else {
			ctx.ReplyTextAndAt("已删除全部任务")
		}
	})
}
