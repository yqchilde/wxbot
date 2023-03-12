package manager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/mid"
	"github.com/yqchilde/wxbot/engine/robot"
)

const (
	RegexOfRemindEveryMonth  = `^设置每月(0?[1-9]|[12][0-9]|3[01])号(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒任务$`
	RegexOfRemindEveryWeek   = `^设置每周(一|二|三|四|五|六|七|日)(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒任务$`
	RegexOfRemindEveryDay    = `^设置每天(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒任务$`
	RegexOfRemindInterval    = `^设置每隔(\d+)(s|秒|m|分|分钟|h|时|d|小时)的提醒任务$`
	RegexOfRemindSpecifyTime = `^设置((20[2-9][0-9]|2100)-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])\s([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒任务$`
	RegexOfRemindExpression  = `^设置表达式\((((\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?)\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?)\)的提醒任务$`
	RegexOfRemindWorkDay     = `^设置工作日(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒任务$`
)

// SetRemindCommand 设置提醒类任务指令
func SetRemindCommand(engine *control.Engine) {
	// 设置每个月的提醒任务
	// Ps: 设置每月8号10:00:00的提醒任务
	engine.OnRegex(RegexOfRemindEveryMonth, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobOfEveryMonth(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置每周的提醒任务
	// 设置每周六20:00:00的提醒任务
	engine.OnRegex(RegexOfRemindEveryWeek, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobOfEveryWeek(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置每天的提醒任务
	// 设置每天10:15:00的提醒任务
	engine.OnRegex(RegexOfRemindEveryDay, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobOfEveryDay(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置每隔多久的提醒任务
	// 设置每隔1小时的提醒任务
	engine.OnRegex(RegexOfRemindInterval, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobForInterval(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置指定时间的一次性提醒任务
	// 设置2023-01-01 15:00:00的提醒任务
	engine.OnRegex(RegexOfRemindSpecifyTime, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobForSpecifyTime(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置自定义cron表达式的提醒任务(6位带秒)
	// 设置表达式(*/10 * * * * *)的提醒任务
	engine.OnRegex(RegexOfRemindExpression, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobForExpression(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 设置工作日的提醒任务
	// 设置工作日10:00:00的提醒任务
	engine.OnRegex(RegexOfRemindWorkDay, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				remind := ctx.MessageString()

				// 设置定时任务
				if _, err := AddCronjobOfEveryDay(ctx, jobTag, matched, func() {
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
						ctx.ReplyTextAndAt(remind)
					}
				}); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:     jobId,
					Tag:    jobTag,
					Type:   JobTypeRemind,
					Desc:   jobDesc,
					WxId:   ctx.Event.FromUniqueID,
					WxName: ctx.Event.FromUniqueName,
					Remind: remind,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, remind))
				job.StartAsync()
				return
			}
		}
	})

	// 删除所有提醒任务
	engine.OnFullMatchGroup([]string{"删除全部提醒任务", "删除所有提醒任务"}, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jobTags []string
		if err := db.Orm.Table("cronjob").Where("wx_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypeRemind).Pluck("tag", &jobTags).Error; err != nil {
			log.Errorf("[CronJob] 删除全部提醒任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部提醒任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("wx_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypeRemind).Delete(&CronJob{}).Error; err != nil {
			log.Errorf("[CronJob] 删除全部提醒任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部提醒任务失败")
			return
		}
		if err := job.RemoveByTagsAny(jobTags...); err != nil {
			log.Errorf("[CronJob] 删除全部提醒任务失败: %v", err)
		} else {
			ctx.ReplyTextAndAt("已删除全部提醒任务")
		}
	})
}
