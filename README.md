# bilibili 直播间 TUI

[关联的bilibili介绍视频](https://www.bilibili.com/video/bv1gG411G7XG)

使用方法 直接下载 releases 中的 bin文件即可

---

风格1: chatroom

![t1](./theme1.png)

风格2: pure

![t2](./theme2.png)

风格3: simple

![t3](./theme3.png)

风格4: info (感谢@soft98-top添加的theme4)

![t4](./theme4.png)

项目文件:

```plaintext
  sender 发送弹幕的实现
  getter 获取弹幕的实现
  ui     TUI的实现
```

使用:

go run main.go

也可以从 参数定义 roomId, theme 优先级高于config(-r roomId, -t theme)

go run main.go -c config.toml -r 9527 -t 1

配置:

默认配置文件: ~/.config/bili/config.toml

参数说明:  
  1. `-c string:configfile`
  2. `-r string:roomId`
  3. `-t int:theme`
  4. `-l int:singleline`
  5. `-s int:showtime`

快捷键:  
  1. \<esc> 退出
  2. <ctrl+c> 退出
  3. <ctrl+u> 清空输入内容
  4. <up> 上一个输入记录
  5. <down> 下一个输入记录

## 类似项目

[zaiic/bili-live-chat](https://github.com/zaiic/bili-live-chat): A bilibili streaming chat tool using TUI written in Rust. 

## 贡献者

- [yaocccc](https://github.com/yaocccc)  
- [soft98-top](https://github.com/soft98-top)
  - [PR#3 增加theme4，修复直播间rank显示](https://github.com/yaocccc/bilibili_live_tui/pull/3)  
- [zaiic](https://github.com/zaiic)
  - [PR#4 更新README，添加类似项目](https://github.com/yaocccc/bilibili_live_tui/pull/4)
- [Ruixi-rebirth](https://github.com/Ruixi-rebirth)
  - [PR#6 自动创建配置文件到 $HOME/.config/bili/config.toml](https://github.com/yaocccc/bilibili_live_tui/pull/6)

## Support: buy me a coffee)

<a href="https://www.buymeacoffee.com/yaocccc" target="_blank">
  <img src="https://github.com/yaocccc/yaocccc/raw/master/qr.png">
</a>
