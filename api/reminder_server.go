package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sejamuchhal/email-reminder/common"
	"github.com/sejamuchhal/email-reminder/storage"
	"github.com/sirupsen/logrus"
)

type ReminderServer struct {
	Storage *storage.Storage
	Logger  *logrus.Entry
}

func (rs *ReminderServer) InitRoutes(r *gin.Engine) {
	r.POST("/reminders/create", rs.CreateReminder)
	r.DELETE("/reminders/delete/:id", rs.DeleteReminder)
	r.GET("/reminders", rs.ListReminders)
}

func StartServer() {
	r := gin.Default()
	conf := common.ConfigureOrDie()
	storage := storage.GetStorageOrDie(conf)
	logger := logrus.New().WithField("component", "reminder_server")
	reminderServer := &ReminderServer{Storage: storage, Logger: logger}
	reminderServer.InitRoutes(r)
	r.Run()
}
