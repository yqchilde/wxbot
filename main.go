package main

import (
	"context"

	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine"
	_ "github.com/yqchilde/wxbot/plugins/baidubaike"   // 百度百科
	_ "github.com/yqchilde/wxbot/plugins/covid19"      // 城市新冠疫情查询
	_ "github.com/yqchilde/wxbot/plugins/crazykfc"     // 肯德基疯狂星期四骚话
	_ "github.com/yqchilde/wxbot/plugins/cronjob"      // 定时任务
	_ "github.com/yqchilde/wxbot/plugins/douyingirl"   // 抖音小姐姐
	_ "github.com/yqchilde/wxbot/plugins/emoticon"     // 表情包原图
	_ "github.com/yqchilde/wxbot/plugins/moyuban"      // 摸鱼办
	_ "github.com/yqchilde/wxbot/plugins/pinyinsuoxie" // 拼音缩写翻译
	_ "github.com/yqchilde/wxbot/plugins/plmm"         // 漂亮妹妹
	_ "github.com/yqchilde/wxbot/plugins/weather"      // 天气查询
)

func main() {
	ctx := context.Background()
	err := engine.Run(ctx, "config.yaml")
	if err != nil {
		log.Fatalf("failed to start robot: %v", err)
	}
}
