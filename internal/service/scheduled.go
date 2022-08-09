package service

import (
	"github.com/yqchilde/pkgs/log"
	"github.com/yqchilde/pkgs/timer"

	"wxBot/internal/config"
	"wxBot/internal/model"
	"wxBot/internal/pkg/holiday"
)

func NewScheduled() {
	scheduled := config.GetScheduled("myb")
	task := timer.NewTimerTask()

	_, err := task.AddTaskByFunc("myb", scheduled[0].Cron, func() {
		if notes, err := holiday.DailyLifeNotes(); err == nil {
			for i := range scheduled[0].Groups {
				model.Groups.SearchByNickName(1, scheduled[0].Groups[i]).SendText(notes)
			}
		}
	})
	if err != nil {
		log.Errorf("NewScheduled add task error: %v", err)
	}
}
