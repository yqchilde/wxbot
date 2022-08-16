package engine

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/yqchilde/pkgs/log"
	"gopkg.in/yaml.v3"

	"github.com/yqchilde/wxbot/engine/config"
)

type FirstConfig config.Config

type Plugin struct {
	context.Context
	context.CancelFunc

	Name      string        // 插件名字
	Version   string        // 插件版本
	Config    config.Plugin // 插件配置
	RawConfig config.Config // 插件原配置
}

func InstallPlugin(conf config.Plugin) *Plugin {
	t := reflect.TypeOf(conf).Elem()
	name := strings.TrimSuffix(t.Name(), "Config")
	plugin := &Plugin{
		Name:   name,
		Config: conf,
	}

	_, pluginFilePath, _, _ := runtime.Caller(1)
	configDir := filepath.Dir(pluginFilePath)
	if parts := strings.Split(configDir, "@"); len(parts) > 1 {
		plugin.Version = parts[len(parts)-1]
	}
	if _, ok := Plugins[name]; ok {
		return nil
	}
	if conf != config.Global {
		Plugins[name] = plugin
		log.Printf("install plugin %v", name)
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
