package cronjob

import (
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/pkgs/timer"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/plugins/moyuban"
)

type CronJob struct {
	engine.PluginMagic
	Enable  bool `yaml:"enable"`
	MoYuBan Job  `yaml:"myb"`
}

type Job struct {
	Cron   string   `yaml:"cron"`
	Groups []string `yaml:"groups"`
}

var (
	pluginInfo = &CronJob{
		PluginMagic: engine.PluginMagic{
			HiddenMenu: true,
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (c *CronJob) OnRegister() {
	// 摸鱼办
	myb := plugin.RawConfig.GetChild("myb")
	{
		cron := myb.Get("cron")
		groups := myb.Get("groups")
		task := timer.NewTimerTask()
		_, err := task.AddTaskByFunc("myb", cron.(string), func() {
			if notes, err := moyuban.DailyLifeNotes("", 0); err == nil {
				for _, val := range groups.([]interface{}) {
					groupList, err := robot.MyRobot.GetGroupList()
					if err != nil {
						panic(err)
					}
					for _, group := range groupList {
						if group.Nickname == val.(string) {
							robot.MyRobot.SendText(group.Wxid, notes)
						}
					}
				}
			}
		})
		if err != nil {
			log.Errorf("NewScheduled add task error: %v", err)
		}
	}
}

func (c *CronJob) OnEvent(msg *robot.Message) {}
