package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

//go:embed fe/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "fe/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})

		router.Run(":8080")

	}()
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("http://127.0.0.1:8080/static/index.html", "", 480, 320, "--remote-allow-origins=*")
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	// cmd := exec.Command(chromePath, "--app=http://127.0.0.1"8080")
	// cmd.Start()
	//处理中断和终止信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//等第一个可以读或可以写的channel进行操作
	select {
	case <-c:
	case <-ui.Done():
	}
	// Wait until UI window is closed
}
