package main

import (
	"github.com/fajarardiyanto/steganography/config"
	"github.com/fajarardiyanto/steganography/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	config.GoogleCloudStorage()
}

func main() {
	r := gin.Default()

	svc := service.New()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/document", svc.RemoveMetadataService)

	if err := r.Run(); err != nil {
		log.Println("unable to connect port :8080")
		return
	}
}
