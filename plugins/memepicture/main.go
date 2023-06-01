package memepicture

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/cryptor"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("memepicture", &control.Options{
		Alias: "表情原图",
		Help: "描述:\n" +
			"表情包不能直接保存下来是个很头疼的事，有时候我们需要保存下来发在其他社交平台，该插件就是这个功能\n\n" +
			"指令:\n" +
			"* 表情原图 -> 30s内发送表情包获取表情原图",
		DataFolder: "memepicture",
	})

	engine.OnFullMatch("表情原图", robot.MustMemePicture).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		imageUrl := ctx.State["image_url"].(string)
		host := ctx.Bot.GetConfig().ServerAddress
		if host == "" {
			ctx.ReplyTextAndAt("请先配置在配置文件中配置serverAddress")
			return
		}

		switch ctx.Bot.GetConfig().Framework.Name {
		case "千寻", "qianxun", "Dean":
			original, thumbnail, err := featureImage(imageUrl, "", engine.GetCacheFolder())
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			thumbnail, err = cryptor.EncryptFilename(ctx.GetFileSecret(), thumbnail)
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			original, err = cryptor.EncryptFilename(ctx.GetFileSecret(), original)
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			thumbnail = fmt.Sprintf("%s/wxbot/static?file=%s", host, thumbnail)
			original = fmt.Sprintf("%s/wxbot/static?file=%s", host, original)
			jumpUrl := fmt.Sprintf("%s/memepicture?img=%s", host, original)
			ctx.ReplyShareLink("快来下载你要的表情原图", "打开后长按图片可保存到本地哦", thumbnail, jumpUrl)
		case "VLW", "vlw":
			original, thumbnail, err := featureImage("", imageUrl, engine.GetCacheFolder())
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			thumbnail, err = cryptor.EncryptFilename(ctx.GetFileSecret(), thumbnail)
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			original, err = cryptor.EncryptFilename(ctx.GetFileSecret(), original)
			if err != nil {
				log.Errorf("获取表情原图失败: %v", err)
				ctx.ReplyTextAndAt("获取表情原图失败")
				return
			}
			thumbnail = fmt.Sprintf("%s/wxbot/static?file=%s", host, thumbnail)
			original = fmt.Sprintf("%s/wxbot/static?file=%s", host, original)
			jumpUrl := fmt.Sprintf("%s/memepicture?img=%s", host, original)
			ctx.ReplyShareLink("快来下载你要的表情原图", "打开后长按图片可保存到本地哦", thumbnail, jumpUrl)
		default:
			ctx.ReplyTextAndAt("暂不支持该框架，请联系管理员")
		}
	})
}

func featureImage(url, b64, cacheDir string) (original, thumbnail string, err error) {
	tmpFile := fmt.Sprintf("%s/tmp_%s%s", cacheDir, time.Now().Local().Format("20060102150405"), ".png")
	if url == "" && b64 == "" {
		return "", "", errors.New("url和b64不能同时为空")
	}
	if url != "" && b64 != "" {
		return "", "", errors.New("url和b64不能同时存在")
	}
	if url != "" {
		// 下载图片
		resp, err := http.Get(url)
		if err != nil {
			return "", "", err
		}

		// 保存图片
		file, err := os.Create(tmpFile)
		if err != nil {
			os.Remove(tmpFile)
			return "", "", err
		}
		if _, err := io.Copy(file, resp.Body); err != nil {
			os.Remove(tmpFile)
			return "", "", err
		}
		resp.Body.Close()
		file.Close()
	}
	if b64 != "" {
		// 保存图片
		if err := utils.Base64ToImage(b64, tmpFile); err != nil {
			os.Remove(tmpFile)
			return "", "", err
		}
	}

	// 检测图片类型
	mime, err := mimetype.DetectFile(tmpFile)
	if err != nil {
		os.Remove(tmpFile)
		return "", "", err
	}
	original = fmt.Sprintf("%s/origin_%s%s", cacheDir, time.Now().Local().Format("20060102150405"), mime.Extension())
	if err := os.Rename(tmpFile, original); err != nil {
		os.Remove(tmpFile)
		return "", "", err
	}

	// 生成缩略图
	if mime.Extension() == ".gif" { // gif
		thumbnail = strings.ReplaceAll(original, "origin", "thumb")
		thumbnail = strings.ReplaceAll(thumbnail, ".gif", ".png")
		return original, thumbnail, gif2Png(original, thumbnail)
	} else {
		thumbnail = original
		return original, thumbnail, nil
	}
}
