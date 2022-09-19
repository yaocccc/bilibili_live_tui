# bilibili 直播间 TUI

![show](show.gif)

[关联的bilibili介绍视频](https://www.bilibili.com/video/bv1gG411G7XG)

项目文件:

```plaintext
  sender 发送弹幕的实现
  getter 获取弹幕的实现
  ui     TUI的实现
```

配置文件: config.toml 仓库内不带 请自己 从 config.toml.demo 复制修改

```plaintext
  Cookie: cookie信息 从web端找一个请求头复制cookie
  RoomId: 指定直播间roomId
```

使用:

go run main.go -c config.toml

快捷键:

1. <esc> 退出
2. <ctrl+c> 退出
3. <ctrl+u> 清空输入内容
4. <up> 上一个输入记录
4. <down> 下一个输入记录
