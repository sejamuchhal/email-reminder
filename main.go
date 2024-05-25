package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sejamuchhal/email-reminder/api"
	"github.com/sejamuchhal/email-reminder/common"
	"github.com/sejamuchhal/email-reminder/storage"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	config := common.ConfigureOrDie()
	storage := storage.GetStorageOrDie(config)
	server := api.ReminderServer{
		Storage: storage,
		Logger:  config.Logger,
	}
	server.InitRoutes(r)
	r.Run()
}
