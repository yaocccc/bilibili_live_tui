# bilibili danmu

服务文件:

```plaintext
  bin/danmu_sender: 发送弹幕的服务
  bin/danmu_geter:  查看弹幕的服务
```

配置文件:

```plaintext
  Cookie: cookie信息 从web端找一个请求头复制cookie
  RoomId: 指定直播间roomId
```

## 使用

需要自己定义config.toml(可以从config.example.toml复制修改)

先执行
./build.sh

再进入 bin 执行

可以参考 里面的run脚本
