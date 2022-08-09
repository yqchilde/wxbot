package service

import (
	"os"

	"wxBot/internal/config"
)

func Init() {
	InitImgDir()
	NewScheduled()
}

func InitImgDir() {
	dir := []string{config.GetPlmmConf().Dir, config.GetEmoticonConf().Dir, config.GetShaoJiConf().Dir}
	for i := range dir {
		err := os.MkdirAll(dir[i], os.ModePerm)
		if err != nil {
			panic("init img dir failed")
		}
	}
}
