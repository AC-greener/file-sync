package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return
}

// 从前端接口的参数里面获取要下载的文件名，然后返给前端
func UploadController(c *gin.Context) {
	if path := c.Params.ByName("path"); path != "" {
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
