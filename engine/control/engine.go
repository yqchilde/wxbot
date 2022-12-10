package control

import (
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/pkgs/utils"

	"github.com/yqchilde/wxbot/engine/robot"
)

type Engine struct {
	en          *robot.Engine // robot engine
	priority    int           // 优先级
	service     string        // 插件服务名
	dataFolder  string        // 数据目录
	cacheFolder string        // 缓存目录
}

var priorityMap = make(map[int]string)
var dataFolderFilter = make(map[string]string)
var cacheFolderFilter = make(map[string]string)

func newEngine(service string, priority int, o *Options[*robot.Ctx]) (e *Engine) {
	s, ok := priorityMap[priority]
	if ok {
		log.Fatal("priority %d is already used by %s", priority, s)
	}
	priorityMap[priority] = service
	log.Debugf("[%s]插件已注册, 优先级: %d", service, priority)

	e = &Engine{
		en:       robot.New(),
		priority: priority,
		service:  service,
	}
	e.en.UsePreHandler(newControl(service, o))

	if o.DataFolder != "" {
		e.dataFolder = "data/plugins/" + o.DataFolder
		if s, ok := dataFolderFilter[e.dataFolder]; ok {
			log.Fatalf("[%s]插件数据目录 %s 已被 %s 占用", service, e.dataFolder, s)
		}
		dataFolderFilter[e.dataFolder] = service
		if err := utils.CheckFolderExists(e.dataFolder); err != nil {
			log.Fatalf("[%s]插件数据目录 %s 创建失败: %v", service, e.dataFolder, err)
		}
	}
	if o.CacheFolder != "" {
		e.cacheFolder = "data/cache/" + o.CacheFolder
		if s, ok := cacheFolderFilter[e.cacheFolder]; ok {
			log.Fatalf("[%s]插件缓存目录 %s 已被 %s 占用", service, e.cacheFolder, s)
		}
		cacheFolderFilter[e.cacheFolder] = service
		if err := utils.CheckFolderExists(e.cacheFolder); err != nil {
			log.Fatalf("[%s]插件缓存目录 %s 创建失败: %v", service, e.cacheFolder, err)
		}
	}
	return
}

// GetDataFolder 获取插件数据目录
func (e *Engine) GetDataFolder() string {
	return e.dataFolder
}

// GetCacheFolder 获取插件缓存目录
func (e *Engine) GetCacheFolder() string {
	return e.cacheFolder
}

// OnMessage 消息触发器
func (e *Engine) OnMessage(rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.On(rules...).SetPriority(e.priority))
}

// OnPrefix 前缀触发器
func (e *Engine) OnPrefix(prefix string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnPrefix(prefix, rules...).SetPriority(e.priority))
}

// OnPrefixGroup 前缀触发器组
func (e *Engine) OnPrefixGroup(prefix []string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnPrefixGroup(prefix, rules...).SetPriority(e.priority))
}

// OnSuffix 后缀触发器
func (e *Engine) OnSuffix(suffix string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnSuffix(suffix, rules...).SetPriority(e.priority))
}

// OnSuffixGroup 后缀触发器组
func (e *Engine) OnSuffixGroup(suffix []string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnSuffixGroup(suffix, rules...).SetPriority(e.priority))
}

// OnCommand 命令触发器
func (e *Engine) OnCommand(commands string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnCommand(commands, rules...).SetPriority(e.priority))
}

// OnCommandGroup 命令触发器组
func (e *Engine) OnCommandGroup(commands []string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnCommandGroup(commands, rules...).SetPriority(e.priority))
}

// OnRegex 正则触发器
func (e *Engine) OnRegex(regexPattern string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnRegex(regexPattern, rules...).SetPriority(e.priority))
}

// OnKeyword 关键词触发器
func (e *Engine) OnKeyword(keyword string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnKeyword(keyword, rules...).SetPriority(e.priority))
}

// OnKeywordGroup 关键词触发器组
func (e *Engine) OnKeywordGroup(keywords []string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnKeywordGroup(keywords, rules...).SetPriority(e.priority))
}

// OnFullMatch 完全匹配触发器
func (e *Engine) OnFullMatch(src string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnFullMatch(src, rules...).SetPriority(e.priority))
}

// OnFullMatchGroup 完全匹配触发器组
func (e *Engine) OnFullMatchGroup(src []string, rules ...robot.Rule) *Matcher {
	return (*Matcher)(e.en.OnFullMatchGroup(src, rules...).SetPriority(e.priority))
}
