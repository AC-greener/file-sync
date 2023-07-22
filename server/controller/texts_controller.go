package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
