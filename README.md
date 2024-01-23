# xbot-ws-server

此插件可以跟 [xbot](https://apifox.com/apidoc/shared-d478def0-67c1-4161-b385-eef8a94e9d17) 进行WS连接，作为WS的服务端，给没有公网IP的xbot发送执行指令。

可以根据自身需求clone后修改ws发送的逻辑，也可以直接下载调用API

# 使用

下载地址 https://github.com/xbot-org/xbot-ws-server/releases

直接运行即可，默认端口`8080`，如需改端口，可以在同目录下新建`.env`文件修改配置

```
PORT=XXXX
```

xbot 的 `.env` 需要配置上 `WS_SERVER` ，对应地址为 `ws://xxxxx:8080/ws`

接口地址为 `http://xxxxx:8080/api`
