package main

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	fmt.Println("Hello, World!")
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
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
		router.POST("/api/v1/texts", TextsController)
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

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw" binding:"required"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		exe, err := os.Executable() //获取当执行路径
		if err != nil {
			log.Fatal(err)
		}

		dir := filepath.Dir(exe) // 获取当前执行文件的目录
		filename := uuid.New().String()
		uploads := filepath.Join(dir, "uploads") // 拼接 uploads 的绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  // 创建 uploads 目录

		if err != nil {
			log.Fatal(err)
		}
		fullpath := filepath.Join("uploads", filename+".txt")                        // 拼接文件的绝对路径（不含 exe 所在目录）
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 将 json.Raw 写入文件

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}
