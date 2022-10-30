package emoticon

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/timer"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type Emoticon struct {
	engine.PluginMagic
	Enable  bool   `yaml:"enable"`
	Dir     string `yaml:"dir"`
	BedUser string `yaml:"bedUser"`
	BedPass string `yaml:"bedPass"`
}

var (
	pluginInfo = &Emoticon{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {表情原图} => 30s内发送表情获取表情原图",
			Commands: []string{"表情原图"},
		},
	}
	plugin      = engine.InstallPlugin(pluginInfo)
	userCommand = make(map[string]string) // 用户指令 key:wx_id val:command
	waitCommand = make(chan *robot.Message)
	mutex       sync.Mutex
)

func (e *Emoticon) OnRegister() {
	conf := plugin.RawConfig
	if err := os.MkdirAll(conf.Get("dir").(string), os.ModePerm); err != nil {
		plugin.Fatalf("init emoticon img dir error: %s", err.Error())
	}

	if err := imageBedLogin(conf.Get("bedUser").(string), conf.Get("bedPass").(string)); err != nil {
		plugin.Fatalf("image bed login error: %s", err.Error())
	}
	imageBedDelete()
}

func (e *Emoticon) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			if addCommand(msg.Content.FromWxid, msg.Content.Msg) {
				return
			}

			if msg.IsSendByPrivateChat() {
				msg.ReplyText("请在30s内发送表情获取表情原图")
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				go waitEmoticon(ctx, cancel, msg)
			} else if msg.IsSendByGroupChat() {
				msg.ReplyTextAndAt("请在30s内发送表情获取表情原图")
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				go waitEmoticon(ctx, cancel, msg)
			}

		}

		if msg.IsEmoticon() {
			for i := range userCommand {
				for j := range pluginInfo.Commands {
					if userCommand[i] == pluginInfo.Commands[j] {
						waitCommand <- msg
						break
					}
				}
			}
		}
	}
}

// 添加用户指令
func addCommand(sender, command string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if val, ok := userCommand[sender]; ok && val == command {
		return true
	} else {
		userCommand[sender] = command
		return false
	}
}

// 移除用户指令
func removeCommand(sender string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(userCommand, sender)
}

func waitEmoticon(ctx context.Context, cancel context.CancelFunc, msg *robot.Message) {
	defer func() {
		cancel()
		removeCommand(msg.Content.FromWxid)
	}()

	for {
		select {
		case <-ctx.Done():
			if msg.IsSendByPrivateChat() {
				msg.ReplyText("30s内未发送表情，获取表情原图失败")
			} else if msg.IsSendByGroupChat() {
				msg.ReplyTextAndAt("30s内未发送表情，获取表情原图失败")
			}
			return
		case msg := <-waitCommand:
			emoticonUrl := msg.Content.Msg[5 : len(msg.Content.Msg)-1]
			b64Str, err := robot.MyRobot.GetFileFoBase64(emoticonUrl)
			if err != nil {
				msg.ReplyText("获取表情原图失败")
				return
			}
			fileName := fmt.Sprintf("%s/emoticon_%s.%s", pluginInfo.Dir, msg.Content.FromWxid, filepath.Ext(emoticonUrl))
			err = base64ToImage(b64Str, fileName)
			if err != nil {
				msg.ReplyText("获取表情原图失败")
				return
			}
			url, err := imageBedUpload(fileName)
			if err != nil {
				msg.ReplyText("获取表情原图失败")
				return
			}
			if err := msg.ReplyShareLink("快来下载你要的表情原图", "打开后长按保存到手机哦", url, url); err != nil {
				imageBedLogin(plugin.RawConfig.Get("bedUser").(string), plugin.RawConfig.Get("bedPass").(string))
				msg.ReplyText(err.Error())
			}
			os.Remove(fileName)
			return
		}
	}
}

var imageBedToken string

func imageBedLogin(username, password string) error {
	type Resp struct {
		Msg   string `json:"msg"`
		Code  int    `json:"code"`
		Token string `json:"token"`
	}
	var resp Resp
	if err := req.C().Post("https://imgbed.link/imgbed/user/login").
		SetFormData(map[string]string{
			"phoneNum": username,
			"pwd":      password,
		}).Do().Into(&resp); err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}
	imageBedToken = resp.Token
	return nil
}

func imageBedUpload(path string) (url string, err error) {
	type Resp struct {
		Total int `json:"total"`
		Rows  []struct {
			Url string `json:"url"`
		} `json:"rows"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var resp Resp
	err = req.C().Post("https://imgbed.link/imgbed/file/upload").
		SetHeader("token", imageBedToken).
		SetFile("file", path).
		Do().Into(&resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(resp.Msg)
	}
	return resp.Rows[0].Url, nil
}

type FileListResp struct {
	Total int `json:"total"`
	Rows  []struct {
		FsId       int    `json:"fsId"`
		UploadTime string `json:"uploadTime"`
	} `json:"rows"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func imageBedDelete() {
	delFunc := func(fsId int) {
		type Resp struct {
			Msg  string `json:"msg"`
			Code int    `json:"code"`
		}
		var resp Resp
		err := req.C().Post("https://imgbed.link/imgbed/file/del").
			SetHeader("token", imageBedToken).
			SetFormData(map[string]string{
				"fsId": strconv.Itoa(fsId),
			}).Do().Into(&resp)
		if err != nil || resp.Code != 0 {
			return
		}
	}

	task := timer.NewTimerTask()
	_, _ = task.AddTaskByFunc("myb", "@daily", func() {
		// 获取文件列表
		client := req.C()
		var fileList FileListResp
		err := client.Post("https://imgbed.link/imgbed/file/mylist").
			SetHeader("token", imageBedToken).
			Do().Into(&fileList)
		if err != nil || fileList.Code != 0 {
			return
		}

		// 删除文件
		for _, file := range fileList.Rows {
			uploadTime, err := time.ParseInLocation("2006-01-02 15:04:05", file.UploadTime, time.Local)
			if err != nil {
				continue
			}
			if time.Now().Local().Sub(uploadTime) > 24*time.Hour {
				delFunc(file.FsId)
			}
		}
	})
}

// base64图片下载到本地方法
func base64ToImage(b64Str, path string) error {
	b, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
