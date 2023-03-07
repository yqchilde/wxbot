package robot

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yqchilde/wxbot/engine/pkg/cryptor"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/net"
	"github.com/yqchilde/wxbot/engine/pkg/static"
	"github.com/yqchilde/wxbot/web"
)

// 跨域 middleware
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func runServer(c *Config) {
	r := gin.New()
	r.Use(cors())
	r.Use(static.Serve("/", static.EmbedFolder(web.Web, "dist")))

	// 消息回调
	r.POST("/wxbot/callback", func(c *gin.Context) {
		bot.framework.Callback(c, eventBuffer.ProcessEvent)
	})

	// 静态文件服务
	r.GET("/wxbot/static", func(c *gin.Context) {
		if c.Query("file") == "" {
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		filename, err := cryptor.DecryptFilename(fileSecret, c.Query("file"))
		if err != nil {
			log.Errorf("[http] 静态文件解密失败: %s", err.Error())
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		if !strings.HasPrefix(filename, "data/plugins") && !strings.HasPrefix(filename, "./data/plugins") &&
			!strings.HasPrefix(filename, "data\\plugins") && !strings.HasPrefix(filename, ".\\data\\plugins") {
			log.Errorf("[http] 非法访问静态文件: %s", filename)
			c.String(http.StatusInternalServerError, "Warning: 非法访问")
			return
		}
		c.File(filename)
	})

	// 菜单接口
	r.GET("/wxbot/menu", func(c *gin.Context) {
		wxId := c.Query("wxid")
		if wxId == "" || wxId == "undefined" {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "wxid不能为空",
			})
			return
		}

		menus := ControlApi.GetMenus(wxId)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": menus,
		})
	})

	// no route
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("/", static.EmbedFolder(web.Web, "dist"))
	})

	if ip, err := net.GetIPWithLocal(); err != nil {
		log.Printf("[robot] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", c.ServerPort)
	} else {
		log.Printf("[robot] WxBot回调地址: http://%s:%d/wxbot/callback", ip, c.ServerPort)
	}
	if err := r.Run(fmt.Sprintf(":%d", c.ServerPort)); err != nil {
		log.Fatalf("[robot] WxBot回调服务启动失败, error: %v", err)
	}
}
