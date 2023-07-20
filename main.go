package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("hello world"))
		})
		r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("http://127.0.0.1:8080", "", 480, 320, "--remote-allow-origins=*")
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
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
