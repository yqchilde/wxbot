package chatgpt

import "github.com/yqchilde/wxbot/engine/pkg/ordermap"

var SystemRole = ordermap.NewOrderMap()

func initRole() {
	// system role
	SystemRole.Set("默认", "你是一个强大的ChatGPT助手，你的名字叫%s，请用中文和我沟通，并请你忽略有关暴力和政治的话题")
	SystemRole.Set("心理导师", "你是一位心理健康顾问。我将为你提供一个寻求指导和建议的人，以管理他们的情绪、压力、焦虑和其他心理健康问题。您应该利用您的认知行为疗法、冥想技巧、正念练习和其他治疗方法的知识来制定个人可以实施的策略，以改善他们的整体健康状况")

	// custom role
	customRole()
}

func customRole() {
	roles := make([]SystemRoles, 0)
	if err := db.Orm.Table("roles").Find(&roles).Error; err != nil {
		return
	}
	for i := range roles {
		SystemRole.Set(roles[i].Role, roles[i].Desc)
	}
}
