package manager

import (
	"fmt"
	"time"

	"github.com/yqchilde/pkgs/timer"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/cronjob"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/mid"
	"github.com/yqchilde/wxbot/engine/robot"
)

var task = timer.NewTimerTask()

type Cronjob struct {
	Id      int64  `gorm:"primary_key"`
	JobId   uint32 `gorm:"column:jid"`
	GroupId string `gorm:"column:gid"`
	Desc    string `gorm:"column:desc"`
	Cron    string `gorm:"column:cron"`
	Remind  string `gorm:"column:remind"`
}

func registerCronjob() {
	engine := control.Register("cronjob", &control.Options[*robot.Ctx]{
		Alias:      "定时任务",
		Help:       "管理员设置定时任务",
		DataFolder: "manager",
	})
	if err := db.Create("cronjob", &Cronjob{}); err != nil {
		log.Fatalf("create cronjob table failed: %v", err)
	}

	go func() {
		c := robot.GetCTX()
		var cronjobs []Cronjob
		if err := db.Orm.Table("cronjob").Find(&cronjobs).Error; err == nil {
			for i := range cronjobs {
				job := cronjobs[i]
				task.AddTaskByFunc("cronjob", job.Cron, func() {
					c.SendText(job.GroupId, job.Remind)
				})
			}
		}
	}()

	engine.OnRegex(`^设置每隔(\d+)(s|秒|m|分|分钟|h|时|d|小时|天)的提醒$`, robot.UserOrGroupAdmin).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		matched := ctx.State["regex_matched"].([]string)
		cronStr := cronjob.ParseToCron(matched[1], matched[2])
		descStr := ctx.MessageString()

		recv, cancel := ctx.EventChannel(ctx.CheckUserSession()).Repeat()
		defer cancel()
		ctx.ReplyText("请问需要提醒什么呢？")
		for {
			select {
			case <-time.After(20 * time.Second):
				ctx.ReplyTextAndAt("操作时间太久了，请重新设置")
				return
			case c := <-recv:
				remind := c.Event.Message.Content
				entryID, err := task.AddTaskByFunc("cronjob", cronStr, func() {
					ctx.ReplyText(remind)
				})
				if err != nil {
					ctx.ReplyText(err.Error())
					return
				}
				db.Orm.Table("cronjob").Create(&Cronjob{
					Id:      mid.UniqueId(),
					JobId:   uint32(entryID),
					GroupId: ctx.Event.FromUniqueID,
					Desc:    descStr,
					Cron:    cronStr,
					Remind:  remind,
				})
				ctx.ReplyTextAndAt(fmt.Sprintf("已为您%s: %s", descStr, remind))
				return
			}
		}
	})

	engine.OnFullMatch("列出所有定时任务").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var cronjobs []Cronjob
		if err := db.Orm.Table("cronjob").Where("gid = ?", ctx.Event.FromUniqueID).Find(&cronjobs).Error; err != nil {
			ctx.ReplyTextAndAt("查询定时任务失败")
			return
		}
		var jobInfo string
		for i := range cronjobs {
			jobInfo += fmt.Sprintf("任务ID: %d\n任务描述: %s\n任务内容: %s\n\n", cronjobs[i].Id, cronjobs[i].Desc, cronjobs[i].Remind)
		}
		if len(cronjobs) == 0 {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务", len(cronjobs)))
		} else {
			ctx.ReplyTextAndAt(fmt.Sprintf("\n当前共有%d个定时任务:\n%s", len(cronjobs), jobInfo))
		}
	})

	engine.OnRegex(`^删除任务 ?(\d+)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		tid := ctx.State["regex_matched"].([]string)[1]
		var jid int
		if err := db.Orm.Table("cronjob").Where("id = ?", tid).Pluck("jid", &jid).Error; err != nil {
			ctx.ReplyTextAndAt("任务ID错误")
			return
		}

		task.RemoveTask("cronjob", jid)
		db.Orm.Table("cronjob").Where("id = ?", tid).Delete(&Cronjob{})
		ctx.ReplyTextAndAt("任务删除成功")
	})

	engine.OnFullMatch("删除全部任务").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var jidList []int
		if err := db.Orm.Table("cronjob").Pluck("jid", &jidList).Error; err != nil {
			ctx.ReplyTextAndAt("删除全部任务失败")
			return
		}
		for i := range jidList {
			task.RemoveTask("cronjob", jidList[i])
		}
		db.Orm.Table("cronjob").Where("gid = ?", ctx.Event.FromUniqueID).Delete(&Cronjob{})
		ctx.ReplyTextAndAt("已删除全部任务")
	})
}
