package engine

import (
	"bytes"
	"context"
	"reflect"

	"github.com/yqchilde/pkgs/log"
	"gopkg.in/yaml.v3"

	"github.com/yqchilde/wxbot/engine/config"
)

type PluginMagic struct {
	Name       string   // 插件名字
	Desc       string   // 插件描述
	Commands   []string // 插件命令
	HiddenMenu bool     // 是否隐藏菜单
}

type Plugin struct {
	context.Context
	context.CancelFunc
	PluginMagic

	Config    config.Plugin // 插件配置
	RawConfig config.Config // 插件原配置
}

func InstallPlugin(conf config.Plugin) *Plugin {
	t := reflect.TypeOf(conf).Elem()
	v := reflect.ValueOf(conf).Elem()

	var p PluginMagic
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Anonymous && t.Field(i).Type.Kind() == reflect.Struct {
			p = v.Field(i).Interface().(PluginMagic)
		}
	}

	plugin := &Plugin{
		PluginMagic: PluginMagic{
			Name:       p.Name,
			Desc:       p.Desc,
			Commands:   p.Commands,
			HiddenMenu: p.HiddenMenu,
		},
		Config: conf,
	}
	if len(plugin.Name) == 0 {
		plugin.Name = t.Name()
	}
	if _, ok := Plugins[plugin.Name]; ok {
		return nil
	}
	if conf != config.Global {
		if !plugin.HiddenMenu && len(plugin.Commands) == 0 {
			log.Errorf("failed to install plugin %s: no commands", plugin.Name)
			return nil
		} else {
			Plugins[plugin.Name] = plugin
			log.Printf("success to install plugin %s", plugin.Name)
		}
	}

	return plugin
}

func (p *Plugin) Assign() {
	p.Context, p.CancelFunc = context.WithCancel(Engine)
	p.RawConfig.Unmarshal(p.Config)

	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(p.Config)
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(&buffer).Decode(&p.RawConfig)
	if err != nil {
		panic(err)
	}
}
