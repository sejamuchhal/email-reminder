package storage

import (
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) *gorm.DB {
	logger := logrus.WithField("method", "MigrateDB")
	DB, err := db.DB()
	if err != nil {
		logger.WithError(err).Fatalln("Could not fetch DB")
	}
	if err = DB.Ping(); err != nil {
		logger.WithError(err).Fatalln("Could not ping DB")
	}
	logger.Info("running migrations")
	options := gormigrate.Options{
		IDColumnName:   "id",
		UseTransaction: true,
	}
	migration := gormigrate.New(db, &options, migrations)
	if err := migration.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration run successfully")
	return db
}

var migrations = []*gormigrate.Migration{
	{
		ID: "202405260450",
		Migrate: func(tx *gorm.DB) error {
			logrus.Println("Creating Reminer table")
			// copy the struct inside the function,
			// so side effects are prevented if the original struct changes during the time
			type Reminder struct {
				Id        int        `gorm:"type:int;primary_key" json:"id"`
				Email     string     `gorm:"size:100;not null" json:"email"`
				Message   string     `gorm:"size:250;not null;"`
				Status    string     `gorm:"size:100;not null;"`
				DueDate   *time.Time `sql:"index"`
				CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
				UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
			}
			if err := tx.Migrator().CreateTable(&Reminder{}); err != nil {
				tx.Rollback()
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("reminders")
		},
	},
}
