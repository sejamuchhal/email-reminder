package storage

import (
	"time"
)

type ReminderStatus string

const (
	StatusCreated ReminderStatus = "created"
	StatusDeleted ReminderStatus = "deleted"
	StatusSent    ReminderStatus = "sent"
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
