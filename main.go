package main

import (
	"github.com/yqchilde/pkgs/config"
	"github.com/yqchilde/pkgs/log"

	model "wxBot/internal/config"
	"wxBot/internal/robot"
	"wxBot/internal/service"
)

func main() {
	// 初始化日志
	config.New(".")
	var conf model.Config
	if err := config.Load("config", &conf); err != nil {
		panic(err)
	}
	model.Conf = &conf
	log.Init(&conf.Log, 2)

	// 初始化service
	service.Init()

	// 初始化机器人
	robot.Init()
}
