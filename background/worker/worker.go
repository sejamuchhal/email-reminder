package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/sejamuchhal/email-reminder/background/services/rabbitmq"
	"github.com/sejamuchhal/email-reminder/common"
	"github.com/sejamuchhal/email-reminder/storage"
	"github.com/streadway/amqp"
)

type Worker struct {
	Storage        *storage.Storage
	EmailSender    *EmailSender
	RabbitMQBroker *rabbitmq.RabbitMQBroker
}

func NewWorker(storage *storage.Storage, emailSender *EmailSender, rmq *rabbitmq.RabbitMQBroker) *Worker {
	return &Worker{
		Storage:        storage,
		EmailSender:    emailSender,
		RabbitMQBroker: rmq,
	}
}

func Run() {
	conf := common.ConfigureOrDie()
	fmt.Printf("Starting background worker with config: %v\n", conf)

	var wg sync.WaitGroup
	var conn *amqp.Connection
    var err error

	logger := conf.Logger
    backoff := time.Millisecond * 500
    maxBackoff := time.Second * 60

    for {
        conn, err = amqp.Dial(conf.RabbitMQURL)
        if err == nil {
            break
        }

        logger.WithError(err).Error("Error while connecting to RabbitMQ, retrying...")
        time.Sleep(backoff)

        backoff *= 2
        if backoff > maxBackoff {
            backoff = maxBackoff
        }
    }

    logger.Info("Successfully connected to RabbitMQ")

	rmq := &rabbitmq.RabbitMQBroker{
		QueueName:  conf.ReminderQueue,
		Connection: conn,
		Logger:     conf.Logger,
	}

	storage := storage.GetStorageOrDie(conf)
	emailSender := NewEmailSender(conf.MailersendAPIKey)

	worker := NewWorker(storage, emailSender, rmq)
	rmq.MsgHandler = worker.ReminderHandler

	wg.Add(1)
	go func() {
		defer wg.Done()
		rmq.Consume()
	}()

	reminderWorker := NewReminderWorker(worker)

	wg.Add(1)
	go func() {
		defer wg.Done()
		reminderWorker.Start()
	}()

	wg.Wait()
}

func newWorker(storage *storage.Storage, emailSender *EmailSender, rmqBroker *rabbitmq.RabbitMQBroker) *Worker {
	return &Worker{
		Storage:        storage,
		EmailSender:    emailSender,
		RabbitMQBroker: rmqBroker,
	}
}
