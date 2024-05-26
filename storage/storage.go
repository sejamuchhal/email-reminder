package storage

import (
	"context"
	"sync"
	"time"

	"github.com/sejamuchhal/email-reminder/common"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var onceConnectAndIndex sync.Once

type Storage struct {
	mainDB *gorm.DB
}

func GetStorageOrDie(config *common.Config) *Storage {
	var storageInstance *Storage
	onceConnectAndIndex.Do(func() {

		db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
		if err != nil {
			logrus.WithError(err).Fatalln("failed to connect to DB")
		}
		if err != nil {
			logrus.WithError(err).Fatalln("failed to connect to DB")
		}
		MigrateDB(db)
		storageInstance = &Storage{mainDB: db}
	})

	return storageInstance
}

func (d *Storage) CreateReminder(ctx context.Context, reminder *Reminder) (*Reminder, error) {
	db := d.mainDB
	err := db.Create(reminder).Error
	return reminder, err
}

func (d *Storage) GetReminderByID(ctx context.Context, ID string) (*Reminder, error) {
	if ID == "" {
		return nil, ErrNoID
	}
	db := d.mainDB
	reminder := &Reminder{}
	err := db.Where("id = ?", ID).Find(reminder).Error
	return reminder, err
}

func (d *Storage) DeleteReminderByID(ctx context.Context, ID string) (int64, error) {
	if ID == "" {
		return 0, ErrNoID
	}
	db := d.mainDB
	reminder := &Reminder{}
	result := db.Model(&Reminder{}).Where("id = ?", ID).Delete(reminder)
	if result.Error != nil {
		return 0, db.Error
	}
	return result.RowsAffected, nil
}

func (d *Storage) ListReminders(ctx context.Context, limit, offset int, status *ReminderStatus) ([]*Reminder, int64, error) {
	db := d.mainDB.Table("reminders")
	var count int64
	reminders := make([]*Reminder, 0, limit)

	if status != nil {
		db = db.Where("status = ?", *status)
	}
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&reminders).Error

	if err != nil {
		return reminders, count, err
	}
	// Setting Limit(-1) and Offset(-1) removes previous limit and offset constraints:
	// https://github.com/go-gorm/gorm/issues/2994
	err = db.Limit(-1).Offset(-1).Count(&count).Error
	return reminders, count, err
}

func (d *Storage) GetRemindersBetween(ctx context.Context, startTime, endTime *time.Time) ([]*Reminder, error) {
	db := d.mainDB
	reminders := make([]*Reminder, 0)

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("status = ? AND due_date BETWEEN ? AND ?", StatusCreated, startTime, endTime).Find(&reminders).Error; err != nil {
			return err
		}

		for _, reminder := range reminders {
			if err := tx.Model(&reminder).Update("status", StatusInProgress).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return reminders, nil
}

func (d *Storage) UpdateReminderStatus(ctx context.Context, ID int, status ReminderStatus) (*Reminder, error) {
	db := d.mainDB
	reminder := &Reminder{}
	err := db.Model(&Reminder{}).Where("id = ?", ID).Update("status", status).First(reminder).Error
	return reminder, err
}
