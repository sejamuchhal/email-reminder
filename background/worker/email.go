package worker

import (
	"context"
	"time"

	"github.com/mailersend/mailersend-go"
)

type EmailSender struct {
	MailerSendServer *mailersend.Mailersend
}

func NewEmailSender(mailersend_api_key string) *EmailSender {
	return &EmailSender{
		MailerSendServer: mailersend.NewMailersend(mailersend_api_key),
	}
}

func (e *EmailSender) SendEmail(to, subject, body string) (*mailersend.Response, error) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  "Seja Muchhal",
		Email: "seja@clarifyme.in",
	}

	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}

	message := e.MailerSendServer.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetText(body)
	message.SetInReplyTo("client-id")

	res, err := e.MailerSendServer.Email.Send(ctx, message)
	if err != nil {
		return nil, err
	}

	return res, nil
}
