package cronjob

import (
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/timer"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/config"
	"github.com/yqchilde/wxbot/engine/robot"
)

type CronJob struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Task   []Task `yaml:"task"`
}

type Task struct {
	Name   string   `yaml:"name"`
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
	task := timer.NewTimerTask()

	tasks := plugin.RawConfig.Get("task").([]interface{})
	for i := range tasks {
		taskConf := tasks[i].(config.Config)
		switch taskConf["name"] {
		case "myb":
			cron := taskConf["cron"].(string)
			groups := taskConf["groups"].([]interface{})
			_, err := task.AddTaskByFunc(taskConf["name"].(string), cron, func() {
				if checkWorkingDay() {
					if url := getMoYuData(); url != "" {
						for i := range groups {
							groupList, err := robot.MyRobot.GetGroupList()
							if err != nil {
								return
							}
							for _, group := range groupList {
								if group.Nickname == groups[i].(string) {
									robot.MyRobot.SendImage(group.Wxid, url)
								}
							}
						}
					}
				}
			})
			if err != nil {
				plugin.Errorf("%s add task error: %v", taskConf["name"].(string), err)
			}
		case "zaoBao":
			cron := taskConf["cron"].(string)
			groups := taskConf["groups"].([]interface{})
			_, err := task.AddTaskByFunc(taskConf["name"].(string), cron, func() {
				if zaoBao, err := getZaoBao(); err == nil {
					for i := range groups {
						groupList, err := robot.MyRobot.GetGroupList()
						if err != nil {
							return
						}
						for _, group := range groupList {
							if group.Nickname == groups[i].(string) {
								robot.MyRobot.SendImage(group.Wxid, zaoBao)
							}
						}
					}
				}
			})
			if err != nil {
				plugin.Errorf("%s add task error: %v", taskConf["name"].(string), err)
			}
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

func getZaoBao() (string, error) {
	type Resp struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		ImageUrl string `json:"imageUrl"`
		Datatime string `json:"datatime"`
	}
	var resp Resp
	if err := req.C().Get("http://dwz.2xb.cn/zaob").Do().Into(&resp); err != nil {
		return "", err
	}
	return resp.ImageUrl, nil
}
