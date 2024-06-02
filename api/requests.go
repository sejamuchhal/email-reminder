package api

type CreateReminderRequest struct {
	Message     string `form:"message" json:"message" binding:"required,min=3,max=100"`
	Email       string `form:"email" json:"email" binding:"required,email"`
	DueDateTime string `form:"due_date_time" json:"due_date_time" binding:"required"`
}

type AuthRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
