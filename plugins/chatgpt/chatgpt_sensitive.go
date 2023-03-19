package chatgpt

import (
	"strings"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var sensitiveWords []string

func initSensitiveWords() {
	sensitiveFile, err := chatGptData.ReadFile("data/sensitive.txt")
	if err != nil {
		log.Errorf("[ChatGPT] 读取敏感词文件失败, error:%s", err.Error())
		return
	}
	// 逐行读取敏感词
	for _, line := range strings.Split(string(sensitiveFile), "\n") {
		if line == "" {
			continue
		}
		sensitiveWords = append(sensitiveWords, line)
	}

	// insert system sensitive words
	for _, word := range sensitiveWords {
		db.Orm.Table("sensitive").FirstOrCreate(&SensitiveWords{Type: 1, Word: word}, "word = ?", word)
	}

	// all sensitive words
	var words []SensitiveWords
	if err := db.Orm.Table("sensitive").Where("deleted = 0").Find(&words).Error; err != nil {
		log.Errorf("[ChatGPT] 获取敏感词失败, error:%s", err.Error())
		return
	}
	sensitiveWords = []string{}
	for _, word := range words {
		sensitiveWords = append(sensitiveWords, word.Word)
	}
}

// 检查敏感词
func checkSensitiveWords(content string) bool {
	for _, word := range sensitiveWords {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

// 将敏感词替换为*
func replaceSensitiveWords(content string) string {
	for _, word := range sensitiveWords {
		if strings.Contains(content, word) {
			content = strings.ReplaceAll(content, word, strings.Repeat("*", len([]rune(word))))
		}
	}
	return content
}

// 设置敏感词相关指令
func setSensitiveCommand(engine *control.Engine) {
	// 查看敏感词列表
	engine.OnRegex("get chatgpt (sensitive|敏感词)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		replyMsg := "当前敏感词列表: " + strings.Join(sensitiveWords, ",")
		log.Debugf("[ChatGPT] 敏感词: %s", replyMsg)
		ctx.ReplyTextAndAt("敏感词无法发出，请查阅日志输出") // 别尝试发出敏感词了，我试了会被吞消息
	})

	// 删除敏感词
	engine.OnRegex("del chatgpt (sensitive|敏感词) (.+)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		word := ctx.State["regex_matched"].([]string)[2]
		find := false
		for i := range sensitiveWords {
			if sensitiveWords[i] == word {
				sensitiveWords = append(sensitiveWords[:i], sensitiveWords[i+1:]...)
				find = true
				break
			}
		}

		if !find {
			ctx.ReplyTextAndAt("敏感词不存在")
			return
		}
		if err := db.Orm.Table("sensitive").Where("word = ?", word).Update("deleted", 1).Error; err != nil {
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		ctx.ReplyTextAndAt("删除敏感词成功")
	})

	// 添加用户自定义敏感词
	engine.OnRegex("set chatgpt (sensitive|敏感词) (.+)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		word := ctx.State["regex_matched"].([]string)[2]
		words := strings.Split(word, ",")
		needs := words
		for i := range words {
			if words[i] == "" {
				needs = append(needs[:i], needs[i+1:]...)
				continue
			}
			for j := range sensitiveWords {
				if sensitiveWords[j] == words[i] {
					needs = append(needs[:i], needs[i+1:]...)
					break
				}
			}
		}

		for i := range needs {
			sensitiveWords = append(sensitiveWords, needs[i])
			db.Orm.Table("sensitive").Where("word = ?", needs[i]).Assign(map[string]interface{}{"deleted": 0}).FirstOrCreate(&SensitiveWords{Type: 2, Word: needs[i]})
		}
		ctx.ReplyTextAndAt("添加敏感词成功")
	})

	// 重置系统敏感词
	engine.OnRegex("reset chatgpt (sensitive|敏感词)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if err := db.Orm.Table("sensitive").Where("type = 1").Delete(&SensitiveWords{}).Error; err != nil {
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		sensitiveWords = []string{}
		initSensitiveWords()
		ctx.ReplyTextAndAt("重置敏感词成功")
	})

	// 删除系统敏感词
	engine.OnRegex("del chatgpt system (sensitive|敏感词)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		tx := db.Orm.Begin()
		if err := tx.Table("sensitive").Where("type = 1").Update("deleted", 1).Error; err != nil {
			tx.Rollback()
			log.Errorf("[ChatGPT] 删除敏感词失败, error:%s", err.Error())
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		var words []SensitiveWords
		if err := tx.Table("sensitive").Where("deleted = 0").Find(&words).Error; err != nil {
			tx.Rollback()
			log.Errorf("[ChatGPT] 删除敏感词失败, error:%s", err.Error())
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		sensitiveWords = []string{}
		for _, word := range words {
			sensitiveWords = append(sensitiveWords, word.Word)
		}
		tx.Commit()
		ctx.ReplyTextAndAt("删除敏感词成功")
	})

	// 删除用户自定义敏感词
	engine.OnRegex("del chatgpt user (sensitive|敏感词)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		tx := db.Orm.Begin()
		if err := tx.Table("sensitive").Where("type = 2").Update("deleted", 1).Error; err != nil {
			tx.Rollback()
			log.Errorf("[ChatGPT] 删除敏感词失败, error:%s", err.Error())
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		var words []SensitiveWords
		if err := tx.Table("sensitive").Where("deleted = 0").Find(&words).Error; err != nil {
			tx.Rollback()
			log.Errorf("[ChatGPT] 删除敏感词失败, error:%s", err.Error())
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		sensitiveWords = []string{}
		for _, word := range words {
			sensitiveWords = append(sensitiveWords, word.Word)
		}
		tx.Commit()
		ctx.ReplyTextAndAt("删除敏感词成功")
	})

	// 删除所有敏感词
	engine.OnRegex("del chatgpt all (sensitive|敏感词)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if err := db.Orm.Table("sensitive").Delete(&SensitiveWords{}).Error; err != nil {
			log.Errorf("[ChatGPT] 删除敏感词失败, error:%s", err.Error())
			ctx.ReplyTextAndAt("删除敏感词失败")
			return
		}
		sensitiveWords = []string{}
		ctx.ReplyTextAndAt("删除敏感词成功")
	})
}
