# UDPXY-Go

UDPXY-Go 是一个使用 Go 语言编写的小型 UDP 代理服务器。它可以获取指定的 IGMP 视频流，并通过 HTTP 提供该流。

## 功能介绍

UDPXY-Go 主要提供以下功能：

1. 接受来自客户端的 HTTP 请求，并解析请求中包含的 UDP 地址。
2. 通过对应的 UDP 地址连接到远程 IGMP 服务器，并获取视频流。
3. 将视频流作为 HTTP 响应返回给客户端。

## 逻辑

客户端通过发送 HTTP GET 请求到 `/udp/<udp-address>` 路径来启动视频流传输。例如，通过访问 `/udp/233.50.201.142:5140` 可以获取 "igmp://233.50.201.142:5140" 对应的视频流。

当服务器接收到这样的请求时，它会建立一个到指定 UDP 地址的 IGMP 连接，并开始读取视频流。然后，服务器会将视频流作为 HTTP 响应的内容发送给客户端。

## 使用方法

首先，你需要编译服务器：

```bash
make build
```

如果你需要为 OpenWrt 构建服务器，可以运行：
```bash
make build-openwrt-amd64 
make build-openwrt-arm 
make build-openwrt-mips
```

你可以在 build/ 目录下找到编译好的二进制文件。

然后，你可以运行服务器：

```bash
./build/udpxy-go
```

# TODO

* HLS
* test RTP