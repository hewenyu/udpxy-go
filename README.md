# udpxy-go
udpxy go 版本的

从 eth0 接收UDP多播流，然后通过 eth1 将这些流以HTTP流的形式提供给另一个网络（192.168.1.0/24）



# 简单设计：

UDPReceiver模块：该模块负责接收UDP多播流。它应该有一个Start方法用于开始接收数据，并将接收到的数据写入一个channel。此外，它还需要能够处理RTP和MPEG-TS负载。

HTTPServer模块：该模块负责处理HTTP请求。它应该有一个Start方法用于开始监听HTTP请求，并从channel中读取数据，将数据作为HTTP响应发送给客户端。此外，它还需要能够解析HTTP命令，并根据命令执行相应的操作。

CommandHandler模块：该模块负责处理HTTP命令。它应该有一个HandleCommand方法用于处理命令，并返回相应的结果。

Logger模块：该模块负责记录日志。它应该有一个Log方法用于记录日志信息。

Main模块：该模块负责启动UDPReceiver和HTTPServer，以及处理它们之间的数据流


# install 

```bash
sudo apt update
# for h264 decoder
sudo apt install libavcodec-dev libavutil-dev libswscale-dev
```