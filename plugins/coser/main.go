package coser

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"modernc.org/mathutil"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
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

	if title, imageUrls := GetCoserInfo(num); title == "" {
		ctx.ReplyTextAndAt("获取 coser 作品失败")
	} else {
		ctx.ReplyTextAndAt(title)
		for _, url := range imageUrls {
			ctx.ReplyImage(url)
			time.Sleep(2 * time.Second)
		}
	}
}

func GetCoserInfo(num int) (string, []string) {
	return getCoserInfo(num, 0)
}

func getCoserInfo(num, retryCount int) (string, []string) {

	// 避免无限递归
	const MaxRetryCount = 3
	if retryCount >= MaxRetryCount {
		log.Error("[coser] 获取失败, 超过最大重试次数")
		return "", nil
	}

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
	// http://pic.yupoo.com 有额度限制，限额后会 302 到一张固定图片，所以此处会重试
	if strings.HasPrefix(data[0], "http://pic.yupoo.com") {
		statusCode := req.C().SetRedirectPolicy(req.NoRedirectPolicy()).SetBaseURL(data[0]).Get().Do().StatusCode
		if statusCode != http.StatusOK {
			return getCoserInfo(num, retryCount+1)
		}
	}

	const MaxCosNum = 10
	min := mathutil.Min(mathutil.Min(num, MaxCosNum), len(data))
	return title, data[0:min]
}
