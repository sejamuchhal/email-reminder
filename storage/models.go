package storage

import (
	"time"
	"gorm.io/gorm"
)

type ReminderStatus string

const (
	StatusCreated    ReminderStatus = "created"
	StatusInProgress ReminderStatus = "inprogress"
	StatusSent       ReminderStatus = "sent"
	StatusFailed     ReminderStatus = "failed"
	StatusDeleted    ReminderStatus = "deleted"
)

type Reminder struct {
	Id        int            `gorm:"type:int;primary_key" json:"id"`
	Email     string         `gorm:"size:100;not null" json:"email"`
	Message   string         `gorm:"size:250;not null;" json:"message"`
	Status    ReminderStatus `gorm:"size:100;not null;" json:"status"`
	DueDate   *time.Time     `sql:"index" json:"due_date"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type User struct {
	gorm.Model
	Email	string `gorm:"unique;size:100;not null" json:"email"`
	Password string `gorm:"size:100;not null" json:"password"`
}