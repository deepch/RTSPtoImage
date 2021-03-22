# RTSPtoImage

RTSP Stream to Image or Mjpeg

use ffmpeg 

![RTSPtoImage image](doc/demo4.png)

## Recommendation

This is an example of getting pictures from rtsp stream

- This is not a working project, this is an example!
  
- For ip cameras, I recommend using the usual get request for the poster. example for dahua [http://<IP address>/onvif/media_service/snapshot?channel=0&subtype=0]

- For this example, I recommend using CUDA, but this will require separate work and will greatly reduce the load on the cpu.

- I recommend limiting the number of frames when using mjpeg or using gpu.

## Installation

### mac os

brew install ffmpeg

### debian / ubuntu

apt install libavcodec-dev libavcodec-ffmpeg56 libavformat-dev  libavformat-ffmpeg56

1.
```bash
go get github.com/deepch/RTSPtoImage
```
2.
```bash
cd $GOPATH/src/github.com/deepch/RTSPtoImage
```
3.
```bash
go run *.go
```
4.
```bash
open web browser http://127.0.0.1:8083
```

## Configuration

### Edit file config.json

format:

```bash
{
  "server": {
    "http_port": ":8083"
  },
  "streams": {
   "H264_AAC": {
      "url": "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov"
    }
  }
}
```

## Limitations

Video Codecs Supported: H264 

Audio Codecs Supported: no

## Team

Deepch - https://github.com/deepch streaming developer

Dmitry - https://github.com/vdalex25 web developer


[![paypal.me/AndreySemochkin](https://ionicabizau.github.io/badges/paypal.svg)](https://www.paypal.me/AndreySemochkin) - You can make one-time donations via PayPal. I'll probably buy a ~~coffee~~ tea. :tea:


