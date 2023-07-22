package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test/server"

	"github.com/zserge/lorca"
)

func main() {
	fmt.Println("Hello, World!")
	go server.Run()

	ui := startChrome()
	// chromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	// cmd := exec.Command(chromePath, "--app=http://127.0.0.1"27149")
	// cmd.Start()
	c := listenToInterupt()
	defer ui.Close()

	//等第一个可以读或可以写的channel进行操作
	select {
	case <-c:
	case <-ui.Done():
	}
}
func startChrome() lorca.UI {
	ui, err := lorca.New("http://127.0.0.1:27149/static/index.html", "", 480, 320, "--remote-allow-origins=*")
	if err != nil {
		log.Fatal(err)
	}
	return ui
}

func listenToInterupt() chan os.Signal {
	//处理中断和终止信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return c
}
