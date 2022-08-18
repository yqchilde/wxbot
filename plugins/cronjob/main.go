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
			Desc:     "🚀 输入 /?? 拼音缩写 => 获取拼音缩写翻译",
			Commands: []string{"/??"},
		},
	}
	plugin = engine.InstallPlugin(pluginInfo)
)

func (c *CronJob) OnRegister(event any) {
	// 摸鱼办
	myb := plugin.RawConfig.GetChild("myb")
	{
		cron := myb.Get("cron")
		groups := myb.Get("groups")
		task := timer.NewTimerTask()
		_, err := task.AddTaskByFunc("myb", cron.(string), func() {
			if notes, err := moyuban.DailyLifeNotes("", 0); err == nil {
				for _, val := range groups.([]interface{}) {
					robot.Groups.SearchByNickName(1, val.(string)).SendText(notes)
				}
			}
		})
		if err != nil {
			log.Errorf("NewScheduled add task error: %v", err)
		}
	}
}

func (c *CronJob) OnEvent(event any) {}
