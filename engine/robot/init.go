package robot

import (
	"errors"
	"io/fs"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/web"
)

var configTemplate = `# 机器人WxId，修改为自己的机器人wxId
botWxId: "你的机器人wxId"
# 机器人名字
botNickname: "Q宝"
# 管理员wxId，多个管理员依次添加，用于管理机器人的wxId
superUsers:
  - "你自己的wxId"
# 管理员命令前缀，匹配系统内置管理员指令需要
commandPrefix: "/"
# 唤醒机器人限制，可选:at，表示群聊必须at机器人才可匹配指令，留空时为默认规则(无特殊要求建议留空)
wakeUpRequire: ""

# 本项目运行时会启动一个HTTP服务，包含一个接收事件回调服务，该配置项为HTTP服务端口
# 在所接入的VX框架中请将回调地址改为 http://[本项目运行服务器IP]:[serverPort]/wxbot/callback
serverPort: 9528
# 本项目运行时会启动一个HTTP服务，包含一个静态图片文件服务，用于将本地图片作为网络图片
# 仅当在插件中使用ctx.ReplyImage(local://[图片路径时])才会用到该项
# 本项目和VX框架运行在同一台服务器时，请填写http://[本机IP]:[serverPort]，否则请填写本项目运行服务器IP，也可以使用域名
serverAddress: "http://192.168.31.12:9528"

# 接入框架配置
framework:
  # 框架选择，可选 千寻、VLW、Dean
  name: "千寻"
  # wxbot主动请求微信框架的地址，比如：http://[运行千寻的服务器ip]:[千寻的HTTP端口]
  apiUrl: "http://192.168.31.8:9527"
  # VX框架HTTP鉴权Token (千寻目前没有，vlw需要)
  apiToken: ""
`

var version string

func init() {
	// 检查配置文件是否存在
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		log.Println("未发现配置文件，已为您生成配置文件，请修改后重新运行程序")
		if err := os.WriteFile("config.yaml", []byte(configTemplate), 0644); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
		os.Exit(0)
	}

	// 打印版本
	log.Printf("当前运行版本: %s", version)

	// 检查web服务
	_, err := web.Web.ReadFile("dist/index.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("web文件夹下的dist文件夹为空，请使用(git clone --recurse-submodules https://github.com/yqchilde/wxbot.git)克隆完整项目")
			return
		}
		log.Fatalf("读取web/dist/index.html失败，可能造成web服务异常: error: %v", err)
	}
	gin.SetMode(gin.ReleaseMode)
}
