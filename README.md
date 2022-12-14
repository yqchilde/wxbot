## What's this?

一个还算有意思的wechat robot项目，诞生于摸鱼 😈

本项目只是用于`Hook`类微信框架处理消息回调，暂时只支持HTTP协议，简而言之`wxbot`
并不能单独实现让您的微信成为Bot，它必须依赖其他`Hook`框架处理回调过来的数据。

* 框架可无侵入扩展，现已集成千寻框架和vlw框架，请参考`framework`目录
* 功能可无侵入扩展，现已集成`plugins`目录下的功能，请参考`plugins`目录

**功能示例：**

![img](https://github.com/yqchilde/wxbot/blob/hook/docs/screenshots.jpg)

<details>
<summary>🎁 已接入框架</summary>

* [x] 千寻框架
* [x] VLW框架

</details>

<details>
<summary>🎁 已有插件</summary>

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
