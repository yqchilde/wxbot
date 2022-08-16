package config

type Engine struct{}

func (cfg *Engine) OnRegister(event any) {}

func (cfg *Engine) OnEvent(event any) {}

var Global = &Engine{}
