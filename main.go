package main

import (
	"context"

	"github.com/yqchilde/wxbot/engine"
	_ "github.com/yqchilde/wxbot/plugins/crazykfc"     // 肯德基疯狂星期四骚话
	_ "github.com/yqchilde/wxbot/plugins/cronjob"      // 漂亮妹妹
	_ "github.com/yqchilde/wxbot/plugins/emoticon"     // 表情包原图
	_ "github.com/yqchilde/wxbot/plugins/jingdong"     // 京东咚咚咚
	_ "github.com/yqchilde/wxbot/plugins/moyuban"      // 摸鱼办
	_ "github.com/yqchilde/wxbot/plugins/pinyinsuoxie" // 拼音缩写翻译
	_ "github.com/yqchilde/wxbot/plugins/plmm"         // 漂亮妹妹
)

func main() {
	ctx := context.Background()
	engine.Run(ctx, "config.yaml")
}
