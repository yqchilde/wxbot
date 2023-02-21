package coser

import (
	"github.com/imroc/req/v3"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
	"math/rand"
	"modernc.org/mathutil"
	"strconv"
	"time"
)

func init() {
	engine := control.Register("coser", &control.Options{
		Alias: "coser",
		Help: "输入 {coser} => 获取 1 张 coser 作品\n" +
			"输入 {coser 5} => 获取 5 张 coser 作品，限制 10 张",
	})

	engine.OnFullMatch("coser").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		reply(ctx, 1)
	})

	engine.OnRegex(`^coser ?(\d+)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if num, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1]); err == nil {
			reply(ctx, num)
		}
	})
}

func reply(ctx *robot.Ctx, num int) {
	if num <= 0 {
		return
	}

	if title, imageUrls := getCoserInfo(num); title == "" {
		ctx.ReplyTextAndAt("获取 coser 作品失败")
	} else {
		ctx.ReplyTextAndAt(title)
		for _, url := range imageUrls {
			ctx.ReplyImage(url)
		}
	}
}

func getCoserInfo(num int) (string, []string) {
	var resp APIResp
	const CoserApiUrl = "http://ovooa.com/API/cosplay/api.php"

	if err := req.C().SetBaseURL(CoserApiUrl).Get().Do().Into(&resp); err != nil {
		log.Errorf("[coser] 请求 %s 失败, err: %v", CoserApiUrl, err)
		return "", nil
	}

	if resp.Code != 1 || len(resp.Data.Data) == 0 {
		log.Errorf("[coser] 请求 %s 失败, resp: %v", CoserApiUrl, resp)
		return "", nil
	}

	title := resp.Data.Title
	data := resp.Data.Data
	if num == 1 {
		return resp.Data.Title, []string{data[rand.Intn(len(resp.Data.Data))]}
	} else {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
		min := mathutil.Min(mathutil.Min(num, 10), len(data))
		return title, data[0:min]
	}
}
