package background

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/sejamuchhal/email-reminder/storage"
)

const WorkerInterval = time.Second * 10

type SendReminderWorker struct {
	Storage     *storage.Storage
	tikcer      *time.Ticker
	stopChan    chan bool
	EmailSender *EmailSender
}

func NewSendReminderWorker(storage *storage.Storage, emailSender *EmailSender) *SendReminderWorker {
	return &SendReminderWorker{
		Storage:     storage,
		tikcer:      time.NewTicker(WorkerInterval),
		stopChan:    make(chan bool),
		EmailSender: emailSender,
	}
}

func (srw *SendReminderWorker) Start() {

	for {
		select {
		case <-srw.tikcer.C:
			srw.SendReminders()
		case <-srw.stopChan:
			srw.tikcer.Stop()
			return
		}
	}
}

func (srw *SendReminderWorker) Stop() {
	srw.stopChan <- true
}

func (srw *SendReminderWorker) SendReminders() {
	logger := log.WithFields(log.Fields{"method": "SendReminders"})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	currentTime := time.Now()
	logger.WithField("currentTime", currentTime).Info("Sending reminders")
	oneIntervalLater := currentTime.Add(WorkerInterval)

	reminders, err := srw.Storage.GetRemindersBetween(ctx, &currentTime, &oneIntervalLater)

	if err != nil {
		logger.WithError(err).Error("failed to get reminders")
		return
	}
	logger.WithField("reminders", len(reminders)).Info("Processing reminders")
	for _, reminder := range reminders {
		err := srw.EmailSender.SendEmail(reminder.Email, "Reminder", reminder.Message)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"reminder_id": reminder.Id,
				"email":       reminder.Email,
			}).WithError(err).Error("Failed to send email reminder")

			if updateErr := srw.updateReminderStatus(ctx, reminder.Id, storage.StatusFailed); updateErr != nil {
				logger.WithFields(logrus.Fields{
					"reminder_id": reminder.Id,
				}).WithError(updateErr).Error("Failed to update reminder status to failed")
			}
			continue
		}

		if updateErr := srw.updateReminderStatus(ctx, reminder.Id, storage.StatusSent); updateErr != nil {
			logger.WithFields(logrus.Fields{
				"reminder_id": reminder.Id,
			}).WithError(updateErr).Error("Failed to update reminder status to sent")
		}
	}

}

func (srw *SendReminderWorker) updateReminderStatus(ctx context.Context, reminderId int, status storage.ReminderStatus) error {
	_, err := srw.Storage.UpdateReminderStatus(ctx, reminderId, status)
	return err
}
