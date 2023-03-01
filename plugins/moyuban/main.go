package moyuban

import (
	"net/http"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("moyu", &control.Options{
		Alias: "摸鱼日历",
		Help: "指令:\n" +
			"* 摸鱼 -> 获取摸鱼办日历\n" +
			"* 摸鱼日历 -> 获取摸鱼办日历",
	})
	engine.OnFullMatchGroup([]string{"摸鱼日历", "摸鱼"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if imageUrl := getMoYuBan(); imageUrl == "" {
			ctx.ReplyTextAndAt("获取摸鱼办日历失败")
		} else {
			ctx.ReplyImage(imageUrl)
		}
	})
}

type MoYu1Resp struct {
	Success bool   `json:"success"`
	Url     string `json:"url"`
}

func getMoYuBan() string {
	resp := req.C().Get("https://api.vvhan.com/api/moyu?type=json").Do()
	if resp.StatusCode != 200 {
		return getMoYuBan2()
	}
	var respData MoYu1Resp
	if err := resp.Into(&respData); err != nil {
		return getMoYuBan2()
	}
	if respData.Url == "" {
		return getMoYuBan2()
	}
	return respData.Url
}

type MoYu2Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MoyuUrl string `json:"moyu_url"`
	} `json:"data"`
}

func getMoYuBan2() string {
	client := req.C()
	resp := client.Get("https://api.j4u.ink/v1/store/other/proxy/remote/moyu.json").Do()
	if resp.StatusCode != 200 {
		return ""
	}
	var respData MoYu2Resp
	if err := resp.Into(&respData); err != nil {
		return ""
	}

	var imageUrl string
	resp = client.SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
		location, err := req.Response.Location()
		if err != nil {
			return err
		}
		imageUrl = location.String()
		return nil
	}).Get(respData.Data.MoyuUrl).Do()
	return imageUrl
}
