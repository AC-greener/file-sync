package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
	c "test/server/controller"
	"test/server/ws"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	hub := ws.NewHub()
	go hub.Run()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.POST("/api/v1/texts", c.TextsController)
	router.GET("/api/v1/addresses", c.AddressController)
	router.GET("/uploads/:path", c.UploadController)
	router.GET("/api/v1/qrcodes", c.QrcodesController)
	router.POST("/api/v1/files", c.FilesController)
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
	router.Run(":27149")
}
