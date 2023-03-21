package choose

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("choose", &control.Options{
		Alias: "选择困难症帮手",
		Help: "指令:\n" +
			"* 帮我选择[选项1]还是[选项2]还是[选项3]还是[选项4]\n" +
			"例:\n" +
			"* 帮我选择可口可乐还是百事可乐\n" +
			"* 帮我选择肯德基还是麦当劳还是必胜客",
	})
	engine.OnPrefix("帮我选择").SetBlock(true).Handle(handle)
}

func handle(ctx *robot.Ctx) {
	rawOptions := strings.Split(ctx.State["args"].(string), "还是")
	if len(rawOptions) == 0 {
		return
	}

	var options = make([]string, 0)
	for count, option := range rawOptions {
		options = append(options, strconv.Itoa(count+1)+". "+option)
	}
	result := rawOptions[rand.Intn(len(rawOptions))]
	err := ctx.ReplyTextAndAt("选项有:\n" + strings.Join(options, "\n") + "\n\n选择结果:\n" + result)
	// 将结果放到匹配队列，触发其它插件
	if err == nil {
		ctx.PushEvent(result)
	}
}
