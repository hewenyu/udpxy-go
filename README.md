# udpxy-go
udpxy go 版本的



# 简单设计：

UDP多播接收器：创建一个goroutine，其任务是监听eth0上的UDP多播地址和端口。当接收到数据时，将其存储在内存中的一个缓冲区（比如一个channel）。

HTTP服务器：创建一个HTTP服务器，监听eth1上特定的端口。当接收到一个HTTP请求时，服务器将从缓冲区中取出数据，并将其作为HTTP响应发送给请求的客户端。
