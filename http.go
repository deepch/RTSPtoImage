package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func serveHTTP() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.LoadHTMLGlob("web/templates/*")
	router.GET("/", func(c *gin.Context) {
		fi, all := Config.list()
		sort.Strings(all)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    fi,
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.GET("/player/:suuid", func(c *gin.Context) {
		_, all := Config.list()
		sort.Strings(all)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    c.Param("suuid"),
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.GET("/play/mjpeg/:suuid", PlayMjpeg)
	router.StaticFS("/static", http.Dir("web/static"))
	err := router.Run(Config.Server.HTTPPort)
	if err != nil {
		log.Fatalln(err)
	}
}
func PlayMjpeg(c *gin.Context) {
	suuid := c.Param("suuid")
	if !Config.ext(suuid) {
		return
	}
	Config.RunIFNotRun(suuid)
	cuuid, ch := Config.clAd(suuid)
	defer Config.clDe(suuid, cuuid)
	c.Header("Content-Type", "multipart/x-mixed-replace;boundary=myboundary")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Connection", "close")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	noVideo := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-noVideo.C:
			log.Println("noVideo")
			return
		case pck := <-ch:
			noVideo.Reset(10 * time.Second)
			header := "\r\n" + "--" + "myboundary" + "\r\n" + "Content-Type: image/jpeg\r\n" + "Content-Length: " + strconv.Itoa(len(*pck)) + "\r\n" + "X-Timestamp: 0.000000\r\n" + "\r\n"
			if n, err := c.Writer.Write([]byte(header)); err != nil || n < 1 {
				return
			}
			if n, err := c.Writer.Write(*pck); err != nil || n < 1 {
				return
			}
		}
	}
}
