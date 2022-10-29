package cronjob

import (
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/pkgs/timer"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
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
			if checkWorkingDay() {
				if url := getMoYuData(); url != "" {
					for _, val := range groups.([]interface{}) {
						groupList, err := robot.MyRobot.GetGroupList()
						if err != nil {
							panic(err)
						}
						for _, group := range groupList {
							if group.Nickname == val.(string) {
								robot.MyRobot.SendImage(group.Wxid, url)
							}
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

func checkWorkingDay() bool {
	type Resp struct {
		Code int `json:"code"`
		Type struct {
			Type int `json:"type"`
		} `json:"type"`
	}
	var resp Resp
	if err := req.C().Get(fmt.Sprintf("http://timor.tech/api/holiday/info?t=%s", time.Now().Format("20060102150405"))).Do().Into(&resp); err != nil {
		return false
	}
	if resp.Code != 0 {
		return false
	}
	if resp.Type.Type == 1 || resp.Type.Type == 2 { // 周末与节假日
		return false
	}
	return true
}

func getMoYuData() (url string) {
	type Resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			MoyuUrl string `json:"moyu_url"`
		} `json:"data"`
	}
	var resp Resp
	if err := req.C().Get("https://api.j4u.ink/v1/store/other/proxy/remote/moyu.json").Do().Into(&resp); err != nil {
		return ""
	}
	if resp.Code != 200 || resp.Message != "success" {
		return ""
	}
	return resp.Data.MoyuUrl
}
