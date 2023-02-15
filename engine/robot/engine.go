package robot

type Engine struct {
	preHandler  []Rule     // 前置处理器
	midHandler  []Rule     // 中间处理器
	postHandler []Handler  // 后置处理器
	block       bool       // 是否阻断后续处理器
	matchers    []*Matcher // 匹配器
}

// New 生成空引擎
func New() *Engine {
	return &Engine{
		preHandler:  []Rule{},
		midHandler:  []Rule{},
		postHandler: []Handler{},
	}
}

var defaultEngine = New()

// UsePreHandler 添加前置处理器，会在 Rule 判断前触发，如果 preHandler 没有通过，则 Rule, Matcher 不会触发
// 可用于分群组管理插件等
func (e *Engine) UsePreHandler(rules ...Rule) {
	e.preHandler = append(e.preHandler, rules...)
}

// UseMidHandler 添加中间处理器，会在 Rule 判断后， Matcher 触发前触发，如果 midHandler 没有通过，则 Matcher 不会触发
// 可用于速率限制等
func (e *Engine) UseMidHandler(rules ...Rule) {
	e.midHandler = append(e.midHandler, rules...)
}

// UsePostHandler 添加后置处理器，会在 Matcher 触发后触发，如果 postHandler 返回 false，则后续的 post handler 不会触发
// 可用于反并发等
func (e *Engine) UsePostHandler(handler ...Handler) {
	e.postHandler = append(e.postHandler, handler...)
}

// SetBlock 设置是否阻断后续处理器
func (e *Engine) SetBlock(block bool) *Engine {
	e.block = block
	return e
}

// On 添加新的匹配器
func On(rules ...Rule) *Matcher { return defaultEngine.On(rules...) }

// On 添加新的匹配器
func (e *Engine) On(rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  rules,
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnMessage 消息触发器
func OnMessage(rules ...Rule) *Matcher { return On(rules...) }

// OnMessage 消息触发器
func (e *Engine) OnMessage(rules ...Rule) *Matcher { return e.On(rules...) }

// OnPrefix 前缀触发器
func OnPrefix(prefix string, rules ...Rule) *Matcher { return defaultEngine.OnPrefix(prefix, rules...) }

// OnPrefix 前缀触发器
func (e *Engine) OnPrefix(prefix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{PrefixRule(prefix)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnPrefixGroup 前缀触发器组
func OnPrefixGroup(prefix []string, rules ...Rule) *Matcher {
	return defaultEngine.OnPrefixGroup(prefix, rules...)
}

// OnPrefixGroup 前缀触发器组
func (e *Engine) OnPrefixGroup(prefix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{PrefixRule(prefix...)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnSuffix 后缀触发器
func OnSuffix(suffix string, rules ...Rule) *Matcher { return defaultEngine.OnSuffix(suffix, rules...) }

// OnSuffix 后缀触发器
func (e *Engine) OnSuffix(suffix string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{SuffixRule(suffix)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnSuffixGroup 后缀触发器组
func OnSuffixGroup(suffix []string, rules ...Rule) *Matcher {
	return defaultEngine.OnSuffixGroup(suffix, rules...)
}

// OnSuffixGroup 后缀触发器组
func (e *Engine) OnSuffixGroup(suffix []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{SuffixRule(suffix...)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCommand 命令触发器
func OnCommand(commands string, rules ...Rule) *Matcher {
	return defaultEngine.OnCommand(commands, rules...)
}

// OnCommand 命令触发器
func (e *Engine) OnCommand(commands string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{CommandRule(commands)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnCommandGroup 命令触发器组
func OnCommandGroup(commands []string, rules ...Rule) *Matcher {
	return defaultEngine.OnCommandGroup(commands, rules...)
}

// OnCommandGroup 命令触发器组
func (e *Engine) OnCommandGroup(commands []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{CommandRule(commands...)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnRegex 正则触发器
func OnRegex(regexPattern string, rules ...Rule) *Matcher {
	return defaultEngine.OnRegex(regexPattern, rules...)
}

// OnRegex 正则触发器
func (e *Engine) OnRegex(regexPattern string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{RegexRule(regexPattern)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnKeyword 关键词触发器
func OnKeyword(keyword string, rules ...Rule) *Matcher {
	return defaultEngine.OnKeyword(keyword, rules...)
}

// OnKeyword 关键词触发器
func (e *Engine) OnKeyword(keyword string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{KeywordRule(keyword)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnKeywordGroup 关键词触发器组
func OnKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	return defaultEngine.OnKeywordGroup(keywords, rules...)
}

// OnKeywordGroup 关键词触发器组
func (e *Engine) OnKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{KeywordRule(keywords...)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnFullMatch 完全匹配触发器
func OnFullMatch(src string, rules ...Rule) *Matcher {
	return defaultEngine.OnFullMatch(src, rules...)
}

// OnFullMatch 完全匹配触发器
func (e *Engine) OnFullMatch(src string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{FullMatchRule(src)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}

// OnFullMatchGroup 完全匹配触发器组
func OnFullMatchGroup(src []string, rules ...Rule) *Matcher {
	return defaultEngine.OnFullMatchGroup(src, rules...)
}

// OnFullMatchGroup 完全匹配触发器组
func (e *Engine) OnFullMatchGroup(src []string, rules ...Rule) *Matcher {
	matcher := &Matcher{
		Engine: e,
		Rules:  append([]Rule{FullMatchRule(src...)}, rules...),
	}
	e.matchers = append(e.matchers, matcher)
	return StoreMatcher(matcher)
}
