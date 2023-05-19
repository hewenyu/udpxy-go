# UDPXY-Go

UDPXY-Go is a small UDP proxy server written in Go language. It can retrieve specified IGMP video streams and serve them over HTTP.

## Features

UDPXY-Go primarily offers the following features:

* Accept HTTP requests from clients and parse the UDP address contained in the requests.
* Connect to the remote IGMP server with the corresponding UDP address and retrieve the video stream.
* Serve the video stream as an HTTP response to the client.

## Logic


The client initiates video stream transmission by sending an HTTP GET request to the `/udp/<udp-address>` path. 

For example, accessing `/udp/233.50.201.142:5140` retrieves the video stream corresponding to `igmp://233.50.201.142:5140`.


## Usage

First, you need to compile the server:

```bash
make build
```

If you need to build the server for OpenWrt, you can run:

```bash
make build-openwrt-amd64 
make build-openwrt-arm 
make build-openwrt-mips
```

You can find the compiled binary files in the `build/` directory.

Then, you can run the server:

```bash
./build/udpxy-go
```

# TODO

* HLS
* test RTP