package background

import (
	"fmt"

	"github.com/sejamuchhal/email-reminder/common"
	"github.com/sejamuchhal/email-reminder/storage"
)

func Run() {

	conf := common.ConfigureOrDie()
	fmt.Printf("Starting background worker with config: %v", conf)
	storage := storage.GetStorageOrDie(conf)
	emailSender := NewEmailSender(conf.MailersendAPIKey)
	worker := NewSendReminderWorker(storage, emailSender)
	worker.Start()
}
