package worker

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

const PushInterval = time.Second * 30

type ReminderWorker struct {
	Ticker   *time.Ticker
	StopChan chan bool
	Worker   *Worker
}

func NewReminderWorker(w *Worker) *ReminderWorker {
	return &ReminderWorker{
		Worker:   w,
		Ticker:   time.NewTicker(PushInterval),
		StopChan: make(chan bool),
	}
}

func (rp *ReminderWorker) Start() {
	for {
		select {
		case <-rp.Ticker.C:
			rp.PushReminders()
		case <-rp.StopChan:
			rp.Ticker.Stop()
			return
		}
	}
}

func (rp *ReminderWorker) Stop() {
	rp.StopChan <- true
}

func (rp *ReminderWorker) PushReminders() {
	logger := log.WithFields(log.Fields{"method": "PushReminders"})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	currentTime := time.Now()
	logger.WithField("currentTime", currentTime).Info("Pushing reminders")
	oneIntervalLater := currentTime.Add(PushInterval)

	reminders, err := rp.Worker.Storage.GetRemindersBetween(ctx, &currentTime, &oneIntervalLater)

	if err != nil {
		logger.WithError(err).Error("failed to get reminders")
		return
	}
	logger.WithField("reminders", len(reminders)).Info("Processing reminders")
	for _, reminder := range reminders {
		reminderJSON, err := json.Marshal(reminder)
		if err != nil {
			logger.WithError(err).Error("failed to marshal reminder")
			continue
		}

		err = rp.Worker.RabbitMQBroker.Publish(reminderJSON)
		if err != nil {
			logger.WithError(err).Error("failed to push reminder")
			continue
		}
		logger.Info("Sent to queue")
	}
}
