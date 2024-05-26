package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sejamuchhal/email-reminder/storage"
)

func (rs *ReminderServer) CreateReminder(c *gin.Context) {
    var reminderReq CreateReminderRequest
    if err := c.ShouldBindJSON(&reminderReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    layout := "January 2, 2006 3:04 PM"
    dueDateTime, err := time.Parse(layout, reminderReq.DueDateTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format"})
        return
    }

    var reminder storage.Reminder
    reminder.Message = reminderReq.Message
    reminder.Email = reminderReq.Email
    reminder.Status = storage.StatusCreated
    reminder.DueDate = &dueDateTime

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    createdReminder, err := rs.Storage.CreateReminder(ctx, &reminder)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Reminder created successfully", "ID": createdReminder.Id})
}

func (rs *ReminderServer) DeleteReminder(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing reminder ID"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rowsAffected, err := rs.Storage.DeleteReminderByID(ctx, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted successfully"})
}

func (rs *ReminderServer) ListReminders(c *gin.Context) {
    rs.Logger.Info("Into ListReminders")
    limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
        return
    }

    offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
        return
    }

    statusParam := c.DefaultQuery("status", "")
    var status *storage.ReminderStatus
    if statusParam != "" {
        status = (*storage.ReminderStatus)(&statusParam)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    reminders, count, err := rs.Storage.ListReminders(ctx, limit, offset, status)
    rs.Logger.WithError(err).Error("error in ListReminders")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"reminders": reminders, "count": count})
}
