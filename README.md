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
<summary>🎁 已接入框架，展开看👇</summary>

* [x] 千寻框架
    * 具体配置查看 `config.yaml` 文件注释说明
    * ![img](https://github.com/yqchilde/wxbot/blob/hook/docs/qianxun.png)
* [x] VLW框架
    * 具体配置查看 `config.yaml` 文件注释说明
    * ![img](https://github.com/yqchilde/wxbot/blob/hook/docs/vlw.png)

</details>

<details>
<summary>🎁 已有插件，展开看 👇</summary>

* [x] [百度百科](https://github.com/yqchilde/wxbot/tree/hook/plugins/baidubaike)
    * 用法：`百度百科 你要查的词`
    * 示例：`百度百科 OCR`
* [x] [ChatGPT聊天](https://github.com/yqchilde/wxbot/tree/hook/plugins/chatgpt)
    * 用法：`# 你要聊的内容`
    * 示例：`# 你好啊`
* [x] [疫情查询](https://github.com/yqchilde/wxbot/tree/hook/plugins/covid19)
    * 用法：`XX疫情`
    * 示例：`济南疫情`
* [x] [KFC疯狂星期四骚话](https://github.com/yqchilde/wxbot/tree/hook/plugins/crazykfc)
    * 用法：`kfc骚话`
    * 示例：`kfc骚话`
* [x] [获取表情原图](https://github.com/yqchilde/wxbot/tree/hook/plugins/memepicture)
    * 用法：输入`表情原图`后30秒内发送表情包(迷因图)
    * 示例：`表情原图`
* [x] [摸鱼办](https://github.com/yqchilde/wxbot/tree/hook/plugins/moyuban)
    * 用法：`摸鱼` `摸鱼办`
    * 用法：`摸鱼办`
* [x] [查拼音缩写](https://github.com/yqchilde/wxbot/tree/hook/plugins/pinyinsuoxie)
    * 用法：`查缩写 你要查的词`
    * 用法：`查缩写 emo`
* [x] [获取美女图片](https://github.com/yqchilde/wxbot/tree/hook/plugins/plmm)
    * 用法：`漂亮妹妹`
    * 示例：`漂亮妹妹`
* [x] [查天气](https://github.com/yqchilde/wxbot/tree/hook/plugins/weather)
    * 用法：`XX天气`
    * 示例：`济南天气`
* [x] [获取每日早报](https://github.com/yqchilde/wxbot/tree/hook/plugins/zaobao)
    * 用法：`早报` `每日早报`
    * 示例：`早报`
* [x] [管理相关](https://github.com/yqchilde/wxbot/tree/hook/plugins/manager)
    * 可以全局设置定时任务
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

  ```bash
  docker run -d \
      --name="wxbot" \
      --restart=always \
      -p 9528:9528 \
      -v $(pwd)/config.yaml:/app/config.yaml \
      -v $(pwd)/data:/app/data \
      yqchilde/wxbot
  ```

## How to develop?

🤔如果您想要扩展自己的插件，可以参考`plugins`目录下的插件

🤔如果您想要扩展其他框架，可以参考`frameworks`目录下的框架

🤔如果您有不想要的插件，可在 `main.go` 上方代码中去掉对应插件的导入(不打算做成动态插件)

## Feature

如果您感觉这个项目有意思，麻烦帮我点一下star  
这个项目待(不)补(完)充(善)很多东西，由于工作关系会抽出时间弄，感谢您发现并使用此仓库

如果您有疑惑可以加交流群讨论，加机器人备注 `wxbot` 邀请您进入

<img src="https://github.com/yqchilde/wxbot/blob/hook/docs/wechat.png" width=30%>

## Thanks

* 非Hook版机器人核心由 [openwechat](https://github.com/eatmoreapple/openwechat) SDK实现，在`nohook`分支，已暂停维护
* Hook版机器人框架我使用的是 ~~《我的框架》已跑路~~，现在用的是千寻，为hook分支

hook分支大量借鉴了一个十分优秀的项目 `ZeroBot-Plugin` 的设计方案 👍🏻，非常感谢，Thanks♪(･ω･)ﾉ
