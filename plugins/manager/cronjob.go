package manager

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/mid"
	"github.com/yqchilde/wxbot/engine/robot"
)

const (
	JobTypeRemind = "remind" // 提醒类任务
	JobTypeFunc   = "func"   // 函数类任务
	JobTypePlugin = "plugin" // 插件类任务

	RegexOfRemindEveryMonth  = `^设置每月(0?[1-9]|[12][0-9]|3[01])号(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒$`
	RegexOfRemindEveryWeek   = `^设置每周(一|二|三|四|五|六|七|日)(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒$`
	RegexOfRemindEveryDay    = `^设置每天(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒$`
	RegexOfRemindInterval    = `^设置每隔(\d+)(s|秒|m|分|分钟|h|时|d|小时)的提醒$`
	RegexOfRemindSpecifyTime = `^设置((20[2-9][0-9]|2100)-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])\s([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])的提醒$`
	RegexOfRemindExpression  = `^设置表达式\((((\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?)\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?\s+(\*(/\d+)?|((\d+(-\d+)?)(,\d+(-\d+)?)*))(/\d+)?)\)的提醒$`
	RegexOfPluginEveryDay    = `^设置每天(([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])执行插件$`
)

type CronJob struct {
	Id      int64  `gorm:"primary_key"`
	Tag     string `gorm:"column:tag"`
	Type    string `gorm:"column:type"`
	Desc    string `gorm:"column:desc"`
	GroupId string `gorm:"column:group_id"`
	Remind  string `gorm:"column:remind"`
	Service string `gorm:"column:service"`
}

func registerCronjob() {
	engine := control.Register("cronjob", &control.Options{
		Alias:      "定时任务",
		Help:       "管理员设置定时任务",
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
		options := control.GetOptionsOnCronjob()
		for i := range cronJobs {
			cronJob := cronJobs[i]
			switch cronJob.Type {
			case JobTypeRemind:
				// 恢复每月的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryMonth).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindOfEveryMonth(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每月提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每周的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryWeek).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindOfEveryWeek(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每周提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复每天的提醒任务
				if matched := regexp.MustCompile(RegexOfRemindEveryDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindOfEveryDay(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复每天提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复间隔提醒任务
				if matched := regexp.MustCompile(RegexOfRemindInterval).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindForInterval(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复间隔提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复指定时间提醒任务
				if matched := regexp.MustCompile(RegexOfRemindSpecifyTime).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindForSpecifyTime(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复指定时间提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}

				// 恢复表达式提醒任务
				if matched := regexp.MustCompile(RegexOfRemindExpression).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddRemindForExpression(ctx, cronJob.Tag, matched, func() {
						ctx.SendText(cronJob.GroupId, cronJob.Remind)
					}); err != nil {
						log.Errorf("恢复表达式提醒任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}
			case JobTypePlugin:
				// 恢复每天的插件任务
				if matched := regexp.MustCompile(RegexOfPluginEveryDay).FindStringSubmatch(cronJob.Desc); matched != nil {
					if _, err := AddPluginOfEveryDay(ctx, cronJob.Tag, matched, func() {
						defer func() {
							if err := recover(); err != nil {
								log.Errorf("执行插件任务失败: %v", string(debug.Stack()))
							}
						}()
						if s, ok := options[cronJob.Service]; ok {
							ctx.Event = &robot.Event{FromUniqueID: cronJob.GroupId}
							s.Options.OnCronjob(ctx)
						}
					}); err != nil {
						log.Errorf("恢复每天插件任务失败: jobId: %d, error: %v", cronJob.Id, err)
					}
				}
			}
		}
		job.StartAsync()
	}()

	// 设置每个月的提醒任务
	// Ps: 设置每月8号10:00:00的提醒
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
				if _, err := AddRemindOfEveryMonth(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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
	// 设置每周六20:00:00的提醒
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
				if _, err := AddRemindOfEveryWeek(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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
	// 设置每天10:15:00的提醒
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
				if _, err := AddRemindOfEveryDay(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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
	// 设置每隔1小时的提醒
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
				if _, err := AddRemindForInterval(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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
	// 设置2023-01-01 15:00:00的提醒
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
				if _, err := AddRemindForSpecifyTime(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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
	// 设置表达式(*/10 * * * * *)的提醒
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
				if _, err := AddRemindForExpression(ctx, jobTag, matched, func() { ctx.ReplyText(remind) }); err != nil {
					ctx.ReplyText(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypeRemind,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Remind:  remind,
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

	// 设置每天的执行插件任务
	// 设置每天08:00:00执行插件
	engine.OnRegex(RegexOfPluginEveryDay, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		jobDesc := ctx.MessageString()
		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		options := control.GetOptionsOnCronjob()
		msg := "请问需要设置哪个插件呢？\n"
		for i := range options {
			msg += options[i].Service + "\n"
		}
		ctx.ReplyText(msg)
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case ctx := <-recv:
				s, ok := options[ctx.MessageString()]
				if !ok {
					ctx.ReplyTextAndAt("没有这个插件服务，请重新设置")
					continue
				}

				jobId := mid.UniqueId()
				jobTag := strconv.Itoa(int(jobId))
				service := ctx.MessageString()

				// 设置定时任务
				if _, err := AddPluginOfEveryDay(ctx, jobTag, matched, func() { s.Options.OnCronjob(ctx) }); err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}

				// 存起来便于服务启动恢复
				if err := db.Orm.Table("cronjob").Create(&CronJob{
					Id:      jobId,
					Tag:     jobTag,
					Type:    JobTypePlugin,
					Desc:    jobDesc,
					GroupId: ctx.Event.FromUniqueID,
					Service: service,
				}).Error; err != nil {
					ctx.ReplyTextAndAt(fmt.Errorf("设置失败: %v", err).Error())
					return
				}
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", jobDesc, service))
				job.StartAsync()
				return
			}
		}
	})

	// 列出当前所有定时任务
	engine.OnFullMatch("列出所有任务").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var cronJobs []CronJob
		if err := db.Orm.Table("cronjob").Where("group_id = ?", ctx.Event.FromUniqueID).Find(&cronJobs).Error; err != nil {
			ctx.ReplyTextAndAt("查询定时任务失败")
			return
		}
		var jobInfo string
		for i := range cronJobs {
			switch cronJobs[i].Type {
			case JobTypeRemind:
				jobInfo += fmt.Sprintf("任务ID: %d\n任务类型: %s\n任务描述: %s\n任务内容: %s\n\n", cronJobs[i].Id, cronJobs[i].Type, cronJobs[i].Desc, cronJobs[i].Remind)
			case JobTypePlugin:
				jobInfo += fmt.Sprintf("任务ID: %d\n任务类型: %s\n任务描述: %s\n任务内容: %s\n\n", cronJobs[i].Id, cronJobs[i].Type, cronJobs[i].Desc, cronJobs[i].Service)
			}
		}
		if len(cronJobs) == 0 {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务", len(cronJobs)))
		} else {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务:\n%s", len(cronJobs), jobInfo))
		}
	})

	// 删除任务 任务ID
	engine.OnRegex(`^删除任务 ?(\d+)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
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

	// 删除所有提醒任务
	engine.OnFullMatchGroup([]string{"删除全部提醒任务", "删除所有提醒任务"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jobTags []string
		if err := db.Orm.Table("cronjob").Where("group_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypeRemind).Pluck("tag", &jobTags).Error; err != nil {
			log.Errorf("[CronJob] 删除全部提醒任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部提醒任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("group_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypeRemind).Delete(&CronJob{}).Error; err != nil {
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

	// 删除所有插件任务
	engine.OnFullMatchGroup([]string{"删除全部插件任务", "删除所有插件任务"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jobTags []string
		if err := db.Orm.Table("cronjob").Where("group_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypePlugin).Pluck("tag", &jobTags).Error; err != nil {
			log.Errorf("[CronJob] 删除全部插件任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部插件任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("group_id = ? AND type = ?", ctx.Event.FromUniqueID, JobTypePlugin).Delete(&CronJob{}).Error; err != nil {
			log.Errorf("[CronJob] 删除全部插件任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部插件任务失败")
			return
		}
		if err := job.RemoveByTagsAny(jobTags...); err != nil {
			log.Errorf("[CronJob] 删除全部插件任务失败: %v", err)
		} else {
			ctx.ReplyTextAndAt("已删除全部插件任务")
		}
	})

	// 删除全部任务
	engine.OnFullMatchGroup([]string{"删除全部任务", "删除所有任务"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jobTags []string
		if err := db.Orm.Table("cronjob").Where("group_id = ?", ctx.Event.FromUniqueID).Pluck("tag", &jobTags).Error; err != nil {
			log.Errorf("[CronJob] 删除全部任务失败: %v", err)
			ctx.ReplyTextAndAt("删除全部任务失败")
			return
		}
		if err := db.Orm.Table("cronjob").Where("group_id = ?", ctx.Event.FromUniqueID).Delete(&CronJob{}).Error; err != nil {
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
