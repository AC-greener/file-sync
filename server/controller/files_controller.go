package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw")
	if err != nil {
		log.Fatal(err)
	}
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
	fullpath := filepath.Join("uploads", filename+filepath.Ext(file.Filename)) // 拼接文件的绝对路径（不含 exe 所在目录）
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})

}
