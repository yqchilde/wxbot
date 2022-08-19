package engine

import (
	"context"
	"os"

	"github.com/yqchilde/pkgs/log"
	"gopkg.in/yaml.v3"

	"github.com/yqchilde/wxbot/engine/config"
)

var (
	// Plugins 所有的插件配置
	Plugins = make(map[string]*Plugin)

	// Engine 复用安装插件逻辑，将全局配置信息注入
	Engine = InstallPlugin(config.Global)
)

func init() {
	log.Default(2)
}

// Run 启动engine，注册plugin
func Run(ctx context.Context, configPath string) (err error) {
	// context bind engine
	Engine.Context = ctx

	// configuration config file
	configRaw, err := os.ReadFile(configPath)
	if err != nil {
		log.Panic("read config file error:", err.Error())
	}
	//var conf config.Config
	conf := new(config.Config)
	if configRaw != nil {
		if err = yaml.Unmarshal(configRaw, conf); err == nil {
			Engine.RawConfig.Unmarshal(config.Global)
		} else {
			log.Panic("parsing yaml error:", err)
		}
	}

	// 合并插件配置
	Engine.RawConfig = config.Struct2Config(config.Global)
	for name, plugin := range Plugins {
		plugin.RawConfig = conf.GetChild(name)
		plugin.Assign()

		if plugin.RawConfig["enable"] != false {
			plugin.Config.OnRegister()
		}
	}

	// 初始化机器人
	InitRobot()

	return
}
