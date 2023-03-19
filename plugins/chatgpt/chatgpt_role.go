package chatgpt

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/yqchilde/wxbot/engine/pkg/ordermap"
	"github.com/yqchilde/wxbot/engine/robot"
)

var SystemRole = ordermap.NewOrderMap()

func initRole() {
	// system role
	SystemRole.Set("默认", "你是一位强大的ChatGPT助手，你的名字叫%s，你将用中文和我沟通，你不会回答任何有关于暴力、色情、政治的话题")
	SystemRole.Set("心理导师", "你是一位心理健康顾问。我将为你提供一个寻求指导和建议的人，以管理他们的情绪、压力、焦虑和其他心理健康问题。您应该利用您的认知行为疗法、冥想技巧、正念练习和其他治疗方法的知识来制定个人可以实施的策略，以改善他们的整体健康状况")

	// custom role
	roles := make([]SystemRoles, 0)
	if err := db.Orm.Table("roles").Find(&roles).Error; err != nil {
		return
	}
	for i := range roles {
		SystemRole.Set(roles[i].Role, roles[i].Desc)
	}
}

// 设置角色相关指令
func setRoleCommand(ctx *robot.Ctx, msg string, command string) {
	switch command {
	case "角色列表":
		replyMsg := "角色列表:\n"
		SystemRole.Each(func(key string, value interface{}) {
			replyMsg += fmt.Sprintf("%s\n", key)
		})
		ctx.ReplyTextAndAt(replyMsg)
	case "当前角色":
		var role string
		if val, ok := chatRoomCtx.Load(ctx.Event.FromUniqueID + "_" + ctx.Event.FromWxId); ok {
			role = val.(ChatRoom).role
		}
		if role == "" {
			ctx.ReplyTextAndAt("当前角色为: 默认")
		} else {
			ctx.ReplyTextAndAt("当前角色为: " + role)
		}
	case "创建角色":
		matched := regexp.MustCompile(`^创建角色\s+(\S+)\s+(.+)$`).FindStringSubmatch(msg)
		role := matched[1]
		if _, ok := SystemRole.Get(role); ok {
			ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]已存在", role))
			return
		}
		desc := matched[2]
		if err := db.Orm.Table("roles").Create(&SystemRoles{Role: role, Desc: desc}).Error; err != nil {
			ctx.ReplyTextAndAt("创建角色失败")
			return
		}
		SystemRole.Set(role, desc)
		ctx.ReplyTextAndAt("创建角色成功")
	case "删除角色":
		matched := regexp.MustCompile(`删除角色\s*(\S+)`).FindStringSubmatch(msg)
		role := matched[1]
		if _, ok := SystemRole.Get(role); !ok {
			ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]不存在", role))
			return
		}
		if err := db.Orm.Table("roles").Where("role = ?", role).Delete(&SystemRoles{}).Error; err != nil {
			ctx.ReplyTextAndAt("删除角色失败")
			return
		}
		SystemRole.Delete(role)
		ctx.ReplyTextAndAt("删除角色成功")
	case "切换角色":
		matched := regexp.MustCompile(`切换角色\s*(\S+)`).FindStringSubmatch(msg)
		role := matched[1]
		if _, ok := SystemRole.Get(role); !ok {
			ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]不存在", role))
			return
		}

		var chatRoom = ChatRoom{
			chatId:   fmt.Sprintf("%s_%s", ctx.Event.FromUniqueID, ctx.Event.FromWxId),
			chatTime: time.Now().Local(),
			role:     role,
			content:  []openai.ChatCompletionMessage{},
		}
		chatRoomCtx.Store(chatRoom.chatId, chatRoom)
		ctx.ReplyTextAndAt("切换角色成功")
	}
}
