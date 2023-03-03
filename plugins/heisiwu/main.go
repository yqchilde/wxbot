package heisiwu

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"modernc.org/mathutil"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	categoryMap = map[string]string{
		"黑丝": "heisi",
		"白丝": "baisi",
		"巨乳": "juru",
		"美足": "meizu",
		"网红": "mcn",
		"jk": "jk",
	}
	categoryKeys = func() []string {
		keys := make([]string, 0, len(categoryMap))
		for k := range categoryMap {
			keys = append(keys, k)
		}
		return keys
	}()
	categoryMatch = strings.Join(categoryKeys, "|")
	categoryRegex = fmt.Sprintf(`^(%s) ?(\d+)$`, categoryMatch)
)

func init() {
	engine := control.Register("heisiwu", &control.Options{
		Alias: "黑丝屋",
		Help: "指令:\n" +
			"* {" + categoryMatch + "} => 获取 1 张作品\n" +
			"* {黑丝 5} => 获取 5 张黑丝作品，限制 10 张\n" +
			"* {巨乳 3} => 获取 3 张巨乳作品，依此类推",
	})

	engine.OnFullMatchGroup(categoryKeys, robot.OnlyPrivate).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		reply(ctx, ctx.State["matched"].(string), 1)
	})

	engine.OnRegex(categoryRegex, robot.OnlyPrivate).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		words := ctx.State["regex_matched"].([]string)
		if num, err := strconv.Atoi(words[2]); err == nil {
			reply(ctx, words[1], num)
		}
	})

	// 启动黑丝屋爬虫
	start()
}

func reply(ctx *robot.Ctx, category string, num int) {
	if num <= 0 || category == "" {
		return
	}

	title, imageUrls := getSetu(categoryMap[category], num)
	if title == "" {
		ctx.ReplyTextAndAt(fmt.Sprintf("获取%s作品失败，请稍后重试", category))
		return
	}

	ctx.ReplyTextAndAt(title)
	for _, url := range imageUrls {
		ctx.ReplyImage(url)
		if dur, err := time.ParseDuration(fmt.Sprintf("%sms", int(rand.Float64()*2000))); err == nil {
			// 等待 2s 内的一个随机数
			time.Sleep(dur)
		}
	}
}

func getSetu(category string, num int) (string, []string) {
	categoryPath := GetPath(StorageFolder, category)
	entries, err := GetSubFolder(categoryPath)

	if err != nil || len(entries) == 0 {
		return "", nil
	}

	title := entries[rand.Intn(len(entries))].Name()
	topicPath := GetPath(categoryPath, title)
	files, err := ReadDir(topicPath)
	if err != nil || len(files) == 0 {
		return "", nil
	}

	sort.Slice(files, func(i, j int) bool {
		left, _ := files[i].Info()
		right, _ := files[j].Info()
		return left.ModTime().Before(right.ModTime())
	})

	const MaxNum = 10
	min := mathutil.Min(mathutil.Min(num, MaxNum), len(files))
	setus := make([]string, 0, min)
	for i := 0; i < min; i++ {
		setus = append(setus, "local://"+GetPath(topicPath, files[i].Name()))
	}
	return title, setus

}
