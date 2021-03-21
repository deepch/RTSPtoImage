package main

import (
	"bytes"
	"errors"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/cgo/ffmpeg"

	"github.com/deepch/vdk/format/rtspv2"
)

var (
	ErrorStreamExitNoVideoOnStream = errors.New("Stream Exit No Video On Stream")
	ErrorStreamExitRtspDisconnect  = errors.New("Stream Exit Rtsp Disconnect")
	ErrorStreamExitNoViewer        = errors.New("Stream Exit On Demand No Viewer")
)

func serveStreams() {
	for k, v := range Config.Streams {
		if !v.OnDemand {
			go RTSPWorkerLoop(k, v.URL, v.OnDemand)
		}
	}
}

func RTSPWorkerLoop(name, url string, OnDemand bool) {
	defer Config.RunUnlock(name)
	for {
		log.Println(name, "Stream Try Connect")
		err := RTSPWorker(name, url, OnDemand)
		if err != nil {
			log.Println(err)
		}
		if OnDemand && !Config.HasViewer(name) {
			log.Println(name, ErrorStreamExitNoViewer)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func RTSPWorker(name, url string, OnDemand bool) error {
	keyTest := time.NewTimer(20 * time.Second)
	clientTest := time.NewTimer(20 * time.Second)
	RTSPClient, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: false, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: false})
	if err != nil {
		return err
	}
	defer RTSPClient.Close()
	if RTSPClient.CodecData != nil {
		Config.coAd(name, RTSPClient.CodecData)
	}
	var VideoStart bool
	var FrameDecoderSingle *ffmpeg.VideoDecoder
	var FrameDecoderStream *ffmpeg.VideoDecoder
	AudioOnly := true

	var videoIDX int
	for i, codec := range RTSPClient.CodecData {
		if codec.Type().IsVideo() {
			AudioOnly = false
		}
		if codec.Type().IsVideo() {
			videoIDX = i
		}
	}

	if !AudioOnly {
		FrameDecoderSingle, err = ffmpeg.NewVideoDecoder(RTSPClient.CodecData[videoIDX].(av.VideoCodecData))
		if err != nil {
			log.Fatalln("FrameDecoderSingle Error", err)
		}

		FrameDecoderStream, err = ffmpeg.NewVideoDecoder(RTSPClient.CodecData[videoIDX].(av.VideoCodecData))
		if err != nil {
			log.Fatalln("FrameDecoderStream Error", err)
		}
	}

	for {
		select {
		case <-clientTest.C:
			if OnDemand && !Config.HasViewer(name) {
				return ErrorStreamExitNoViewer
			}
		case <-keyTest.C:
			return ErrorStreamExitNoVideoOnStream
		case signals := <-RTSPClient.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				Config.coAd(name, RTSPClient.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return ErrorStreamExitRtspDisconnect
			}
		case packetAV := <-RTSPClient.OutgoingPacketQueue:
			if AudioOnly || packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
			}
			if packetAV.IsKeyFrame {
				VideoStart = true
			}
			//sample single frame decode encode to jpeg save on disk //
			if packetAV.IsKeyFrame {
				if pic, err := FrameDecoderSingle.DecodeSingle(packetAV.Data); err == nil && pic != nil {
					if out, err := os.Create("./output.jpg"); err == nil {
						if err = jpeg.Encode(out, &pic.Image, nil); err == nil {
							log.Println("Poster Saved On Disk Ready ./output.jpg")
						}
					}
				}
			}
			//sample stream video decode encode jpeg and play mjpeg over http //
			if VideoStart {
				if packetAV.Idx == int8(videoIDX) {
					if pic, err := FrameDecoderStream.DecodeSingle(packetAV.Data); err == nil && pic != nil {
						buf := new(bytes.Buffer)
						if err = jpeg.Encode(buf, &pic.Image, nil); err == nil {
							Config.cast(name, buf.Bytes())
						}
					}
				}
			}
		}
	}
}
