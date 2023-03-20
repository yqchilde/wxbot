package wordcloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("wordcloud", &control.Options{
		Alias: "词云",
		Help: "描述:\n" +
			"想了解一下当前聊天室的聊天热点？快来试一试\n\n" +
			"指令:\n" +
			"* 热词 -> 获取当前聊天室热词，默认当前聊天室Top30条\n" +
			"* 热词 top [10] -> 获取当前聊天室热词前10条\n" +
			"* 热词 id [xxx] -> 获取指定聊天室热词\n" +
			"* 热词 id [xxx] top [10] -> 获取指定聊天室热词前10条\n" +
			"* 引用群里某个人消息并回复 热词，获取该用户在群里的热词",
		DataFolder: "wordcloud",
	})

	engine.OnRegex(`^热词(?:\s+id\s+(\S+))?(?:\s+top\s+(\d+))?$|^热词\s+top\s+(\d+)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		id := ctx.State["regex_matched"].([]string)[1]
		top, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)

		// todo 5分钟允许拿一次数据就行了，有空在搞

		if id == "" && top == 0 { // 获取当前群，top默认30
			id = ctx.Event.FromUniqueID
			top = 30
		} else if id != "" && top == 0 { // 获取指定群，top默认30
			top = 30
		} else if id == "" && top != 0 { // 获取当前群，top指定
			id = ctx.Event.FromUniqueID
		} else if id != "" && top != 0 { // 获取指定群，top指定
			// do nothing
		}

		// 获取历史记录-文本消息
		var record []robot.MessageRecord
		if ctx.IsReference() { // 查询聊天室指定用户
			r, err := ctx.GetRecordHistory(&robot.RecordConditions{
				FromWxId:   id,
				SenderWxId: ctx.Event.ReferenceMessage.ChatUser,
				CreatedAt:  time.Now().Local().Format("2006-01-02"),
			})
			if err != nil {
				log.Errorf("获取[%s]热词失败: %v", id, err)
				ctx.ReplyText("获取热词失败")
				return
			}
			record = r
		} else { // 查询聊天室所有消息
			r, err := ctx.GetRecordHistory(&robot.RecordConditions{
				FromWxId:  id,
				CreatedAt: time.Now().Local().Format("2006-01-02"),
			})
			if err != nil {
				log.Errorf("获取[%s]热词失败: %v", id, err)
				ctx.ReplyText("获取热词失败")
				return
			}
			record = r
		}

		// 整理文本消息
		var words string
		for _, msg := range record {
			// 剔除消息中的表情
			for _, emoji := range robot.EmojiSymbol {
				msg.Content = strings.ReplaceAll(msg.Content, emoji, "")
			}
			// 剔除消息中的艾特
			if strings.HasPrefix(msg.Content, "@") {
				if strings.Contains(msg.Content, " ") {
					msg.Content = msg.Content[strings.Index(msg.Content, " "):]
				} else {
					msg.Content = msg.Content[1:]
				}
			}
			words += msg.Content + " "
		}

		// 获取热词图
		resp := req.C().Post("https://bot.yqqy.top/api/wordcloud").SetBody(map[string]interface{}{"words": words, "count": top}).Do()
		if resp.GetStatusCode() != 200 {
			log.Errorf("获取[%s]热词失败，HTTP Status: %v", id, resp.GetStatusCode())
			ctx.ReplyText("获取热词失败")
			return
		}
		if code := gjson.Get(resp.String(), "code").Int(); code != 200 {
			if code == 401 {
				ctx.ReplyText("数据太少，再聊会儿吧")
				return
			}
			log.Errorf("获取[%s]热词失败: code: %v", id, code)
			ctx.ReplyText("获取热词失败")
			return
		}

		// 保存图片
		imgB64 := gjson.Get(resp.String(), "data.image").String()
		filename := fmt.Sprintf("%s/%s_%s.png", engine.GetCacheFolder(), ctx.Event.FromUniqueID, time.Now().Local().Format("20060102"))
		if err := utils.Base64ToImage(imgB64, filename); err != nil {
			log.Errorf("保存图片失败: %v", err)
			ctx.ReplyText("获取热词失败")
			return
		}

		// 发送图片
		if ctx.IsReference() {
			ctx.ReplyText(fmt.Sprintf("[%s]今天的热词如图", ctx.Event.ReferenceMessage.DisplayName))
		}
		ctx.ReplyImage("local://" + filename)
	})
}
