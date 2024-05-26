package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sejamuchhal/email-reminder/storage"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (w *Worker) ReminderHandler(queue string, msg amqp.Delivery, err error) {
	logger := log.WithFields(log.Fields{"method": "ReminderHandler"})

	if err != nil {
		logger.WithError(err).Error("Error occurred in RMQ consumer")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger.Infof("Message received on '%s' queue: %s", queue, string(msg.Body))
	var reminder storage.Reminder
	err = json.Unmarshal(msg.Body, &reminder)
	if err != nil {
		logger.WithError(err).Error("Error while unmarshalling reminder")
		return
	}

	_, err = w.EmailSender.SendEmail(reminder.Email, "Reminder", reminder.Message)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"reminder_id": reminder.Id,
			"email":       reminder.Email,
		}).WithError(err).Error("Failed to send email reminder")

		if updateErr := w.updateReminderStatus(ctx, reminder.Id, storage.StatusFailed); updateErr != nil {
			logger.WithFields(logrus.Fields{
				"reminder_id": reminder.Id,
			}).WithError(updateErr).Error("Failed to update reminder status to failed")
		}
	}

	if updateErr := w.updateReminderStatus(ctx, reminder.Id, storage.StatusSent); updateErr != nil {
		logger.WithFields(logrus.Fields{
			"reminder_id": reminder.Id,
		}).WithError(updateErr).Error("Failed to update reminder status to sent")
	}
	logger.WithFields(logrus.Fields{"reminder_id": reminder.Id, "email": reminder.Email}).Info("Email sent successfully")

}

func (w *Worker) updateReminderStatus(ctx context.Context, reminderId int, status storage.ReminderStatus) error {
	_, err := w.Storage.UpdateReminderStatus(ctx, reminderId, status)
	return err
}
