package config

import "github.com/yqchilde/pkgs/log"

var Conf *Config

type Config struct {
	Emoticon  Emoticon    `yaml:"emoticon"`
	Plmm      Plmm        `json:"plmm"`
	ShaoJi    ShaoJi      `json:"shaoji"`
	Log       log.Config  `yaml:"log"`
	Scheduled []Scheduled `yaml:"scheduled"`
}

type Emoticon struct {
	Dir string `yaml:"dir"`
}

type Plmm struct {
	Dir       string `yaml:"dir"`
	Url       string `yaml:"url"`
	AppId     string `yaml:"appId"`
	AppSecret string `yaml:"appSecret"`
}

type ShaoJi struct {
	Dir string `yaml:"dir"`
	Url string `yaml:"url"`
}

type Scheduled struct {
	Name   string   `yaml:"name"`
	Cron   string   `yaml:"cron"`
	Groups []string `yaml:"groups"`
}

func GetEmoticonConf() *Emoticon {
	return &Conf.Emoticon
}

func GetPlmmConf() *Plmm {
	return &Conf.Plmm
}

func GetShaoJiConf() *ShaoJi {
	return &Conf.ShaoJi
}

func GetScheduled(name ...string) []Scheduled {
	if len(name) == 0 {
		return Conf.Scheduled
	} else {
		var scheduled []Scheduled
		for _, v := range Conf.Scheduled {
			if v.Name == name[0] {
				scheduled = append(scheduled, v)
			}
		}
		return scheduled
	}
}
