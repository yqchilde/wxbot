package config

import (
	"github.com/yqchilde/wxbot/engine/robot"
)

type Engine struct{}

func (cfg *Engine) OnRegister() {}

func (cfg *Engine) OnEvent(msg *robot.Message) {}

var Global = &Engine{}
