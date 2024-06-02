package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sejamuchhal/email-reminder/storage"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (rs *ReminderServer) CreateReminder(c *gin.Context) {
	var reminderReq CreateReminderRequest
	if err := c.ShouldBindJSON(&reminderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, err := time.Parse(layout, reminderReq.DueDateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format, please try again with format: January 2, 2006 3:04 PM MST"})
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

	c.JSON(http.StatusCreated, gin.H{"message": "Reminder created successfully", "created_reminder": createdReminder})
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reminders": reminders, "count": count})
}

func (rs *ReminderServer) SignUp(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		rs.Logger.WithError(err).Error("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	user := storage.User{
		Email:    req.Email,
		Password: string(hash),
	}

	res, err := rs.Storage.CreateUser(context.Background(), &user)
	if err != nil {
		rs.Logger.WithError(err).Error("Failed to create user")
		if err == storage.ErrUserExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with this email"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user_id": res.ID})
}

func (rs *ReminderServer) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := rs.Storage.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		}
		rs.Logger.WithError(err).Error("Error while fetching user by email")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(rs.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
