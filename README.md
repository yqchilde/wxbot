## What's this?

一个还算有意思的wechat robot项目，理想将其设计成一个多功能的机器人助手 😈

* 🤨框架可无侵入扩展，现已接入`千寻框架`和`vlw框架`，请参考`framework`目录
* 🤨功能可无侵入扩展，现已集成`plugins`目录下的功能，请参考`plugins`目录

**🔔 注意：**

1. 使用本项目之前需要您已经配置好相关的 `微信的hook` 类软件，那么只需要在这类软件上设置回调地址即可
2. 本项目已接入`vlw`、`千寻`两个框架，如果您有其他框架，可自行添加(参考`framework`目录，实现`IFramework`接口即可)，或联系我添加
3. 本项目不提供任何`hook`类软件，您需要利用搜索引擎自行寻找
4. 本项目暂时只支持HTTP协议，关于websocket协议支持目前不考虑
5. 简而言之，本项目是一个消息处理的中间件，微信消息监听获取是从框架获取
6. 本项目仅供学习交流使用，不得用于商业用途，否则后果自负
7. 使用本项目造成封禁账号等后果（项目立项到现在，作者还没出现过异常），本项目不承担任何责任，实际上您使用任何非官方的微信机器人都有可能造成账号封禁，所以请谨慎使用
8. 如果您阅读了上面的内容，觉得没有问题，那么请继续阅读下面的内容

**功能示例：**

![img](https://github.com/yqchilde/wxbot/blob/hook/docs/screenshots.jpg)

<details>
<summary>🎁 已对接API，展开看👇</summary>

```go
type IFramework interface {
	// Callback 这是消息回调方法，vx框架回调消息转发给该Server
	Callback(func(*Event, IFramework))

	// GetMemePictures 获取表情包图片地址(迷因图)
	// return: 图片链接(网络URL或图片base64)
	GetMemePictures(message *Message) string

	// SendText 发送文本消息
	// toWxId: 好友ID/群ID
	// text: 文本内容
	SendText(toWxId, text string) error

	// SendTextAndAt 发送文本消息并@，只有群聊有效
	// toGroupWxId: 群ID
	// toWxId: 好友ID/群ID/all
	// toWxName: 好友昵称/群昵称，留空为自动获取
	// text: 文本内容
	SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error

	// SendImage 发送图片消息
	// toWxId: 好友ID/群ID
	// path: 图片路径
	SendImage(toWxId, path string) error

	// SendShareLink 发送分享链接消息
	// toWxId: 好友ID/群ID
	// title: 标题
	// desc: 描述
	// imageUrl: 图片链接
	// jumpUrl: 跳转链接
	SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error

	// SendFile 发送文件消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地文件绝对路径
	SendFile(toWxId, path string) error

	// SendVideo 发送视频消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地视频文件绝对路径
	SendVideo(toWxId, path string) error

	// SendEmoji 发送表情消息
	// toWxId: 好友ID/群ID/公众号ID
	// path: 本地动态表情文件绝对路径
	SendEmoji(toWxId, path string) error

	// SendMusic 发送音乐消息
	// toWxId: 好友ID/群ID/公众号ID
	// name: 音乐名称
	// author: 音乐作者
	// app: 音乐来源(VLW需留空)，酷狗/wx79f2c4418704b4f8，网易云/wx8dd6ecd81906fd84，QQ音乐/wx5aa333606550dfd5
	// jumpUrl: 音乐跳转链接
	// musicUrl: 网络歌曲直链
	// coverUrl: 封面图片链接
	SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error

	// SendMiniProgram 发送小程序消息
	// toWxId: 好友ID/群ID/公众号ID
	// ghId: 小程序ID
	// title: 标题
	// content: 内容
	// imagePath: 图片路径, 本地图片路径或网络图片URL
	// jumpPath: 小程序点击跳转地址，例如：pages/index/index.html
	SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error

	// SendMessageRecord 发送消息记录
	// toWxId: 好友ID/群ID/公众号ID
	// title: 仅供电脑上显示用，手机上的话微信会根据[显示昵称]来自动生成 谁和谁的聊天记录
	// dataList:
	// 	- wxid: 发送此条消息的人的wxid
	// 	- nickName: 显示的昵称(可随意伪造)
	// 	- timestamp: 10位时间戳
	// 	- msg: 消息内容
	SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error

	// SendMessageRecordXML 发送消息记录(XML方式)
	// toWxId: 好友ID/群ID/公众号ID
	// xmlStr: 消息记录XML代码
	SendMessageRecordXML(toWxId, xmlStr string) error

	// SendFavorites 发送收藏消息
	// toWxId: 好友ID/群ID/公众号ID
	// favoritesId: 收藏夹ID
	SendFavorites(toWxId, favoritesId string) error

	// SendXML 发送XML消息
	// toWxId: 好友ID/群ID/公众号ID
	// xmlStr: XML代码
	SendXML(toWxId, xmlStr string) error

	// SendBusinessCard 发送名片消息
	// toWxId: 好友ID/群ID/公众号ID
	// targetWxId: 目标用户ID
	SendBusinessCard(toWxId, targetWxId string) error

	// AgreeFriendVerify 同意好友验证
	// v3: 验证V3
	// v4: 验证V4
	// scene: 验证场景
	AgreeFriendVerify(v3, v4, scene string) error

	// InviteIntoGroup 邀请好友加入群组
	// groupWxId: 群ID
	// wxId: 好友ID
	// typ: 邀请类型，1-直接拉，2-发送邀请链接
	InviteIntoGroup(groupWxId, wxId string, typ int) error

	// GetObjectInfo 获取对象信息
	// wxId: 好友ID/群ID/公众号ID
	// return: ObjectInfo, error
	GetObjectInfo(wxId string) (*ObjectInfo, error)
}
```

</details>

<details>
<summary>🎁 已接入框架，展开看👇</summary>

* [x] 千寻框架
    * 具体配置查看 `config.yaml` 文件注释说明
    * ![img](https://github.com/yqchilde/wxbot/blob/hook/docs/qianxun.png)
* [x] VLW框架
    * 具体配置查看 `config.yaml` 文件注释说明
    * ![img](https://github.com/yqchilde/wxbot/blob/hook/docs/vlw.png)

</details>

<details open>
<summary>🎁 已有插件 👇</summary>

* [x] [百度百科-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/baidubaike)
    * 用法：发送`百度百科 你要查的词`，例如：`百度百科 OCR`
* [x] [ChatGPT聊天-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/chatgpt)
    * 用法：发送`开始ChatGPT会话`，然后就可以和机器人连续对话聊天了
* [x] [KFC疯狂星期四骚话-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/crazykfc)
    * 用法：发送`kfc骚话`，获取一条v50骚话
* [x] [获取表情原图-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/memepicture)
    * 用法：发送`表情原图`后30秒内发送一张表情包(迷因图)，即可获取原图
* [x] [摸鱼办-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/moyuban)
    * 用法：发送`摸鱼`或`摸鱼办`，即可获取一张摸鱼办图片
* [x] [查拼音缩写-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/pinyinsuoxie)
    * 用法：发送`查缩写 你要查的词`，即可获取拼音缩写含义
* [x] [获取美女图片-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/plmm)
    * 用法：发送`漂亮妹妹`，即可获取一张美女图片
* [x] [查天气-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/weather)
    * 用法：发送`XX天气`，即可获取XX地区的天气情况，例如：`济南天气`
* [x] [获取每日早报-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/zaobao)
    * 用法：发送`早报`或`每日早报`，即可获取每日早报
* [x] [管理相关-点击查看详情](https://github.com/yqchilde/wxbot/tree/hook/plugins/manager)
    * 可以全局设置定时任务
        * 设置每月8号10:00:00的提醒 
        * 设置每周三10:00:00的提醒 
        * 设置每天10:00:00的提醒 
        * 设置每隔1小时的提醒 
    * 可以全局监听好友添加邀请拉群等

</details>

## How to use?

### 本地运行

1. 拷贝代码

    ```bash
    git clone https://github.com/yqchilde/wxbot.git
    ```

2. 配置`config.yaml`

3. `go run main.go` 或自行build

### Docker运行

1. 一键脚本启动

```shell
bash -c "$(curl -fsSL https://raw.fastgit.org/yqchilde/wxbot/hook/docker/run.sh)"
```

2. 命令启动，注意提前配置`config.yaml`,否则会报错
  ```shell
  docker run -d \
      --name="wxbot" \
      -p 9528:9528 \
      -v $(pwd)/config.yaml:/app/config.yaml \
      -v $(pwd)/data:/app/data \
      yqchilde/wxbot:latest
  ```

## How to develop?

🤔如果您想要扩展自己的插件，可以参考`plugins`目录下的插件

🤔如果您想要扩展其他框架，可以参考`frameworks`目录下的框架

🤔如果您有不想要的插件，可在 `main.go` 上方代码中去掉对应插件的导入(不打算做成动态插件)

### 调试-环境变量

| 环境变量名 | 变量类型 | 说明                                                         |
| ---------- | -------- | ------------------------------------------------------------ |
| DEBUG      | bool     | 优先级大于其他`DEBUG_`开头的变量，开启后开启所有DEBUG模式<br />用于调试HTTP请求和调用日志文件名和行号 |
| DEBUG_LOG  | bool     | 用于调试调用日志文件名和行号                                 |

## Feature

如果您感觉这个项目有意思，麻烦帮我点一下star  
这个项目待(不)补(完)充(善)很多东西，由于工作关系会抽出时间弄，感谢您发现并使用此仓库

如果您有疑惑可以加Q群讨论

<img src="https://github.com/yqchilde/wxbot/blob/hook/docs/qq.jpg" width=30%>

## Thanks

* 非Hook版机器人核心由 [openwechat](https://github.com/eatmoreapple/openwechat) SDK实现，在`nohook`分支，已暂停维护
* Hook版机器人框架我使用的是 ~~《我的框架》已跑路~~，现在用的是千寻，为hook分支

hook分支大量借鉴了一个十分优秀的项目`ZeroBot-Plugin`的设计方案 👍🏻，其中很多基础代码来自`ZeroBot-Plugin`，在此基础上扩展了支持`wechat`的方式，非常感谢，Thanks♪(･ω･)ﾉ
