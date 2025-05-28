package email

import (
	"context"
	"github.com/ZekebayevYe/notification-service/internal/app"
	"github.com/mailersend/mailersend-go"
	"log"
	"os"
	"time"
)

type Mailer struct {
	client *mailersend.Mailersend
}

func NewMailer() *Mailer {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	return &Mailer{client: mailersend.NewMailersend(apiKey)}
}

func (m *Mailer) SendNotification(n app.Notification, to []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := m.client.Email.NewMessage()
	msg.SetFrom(mailersend.From{
		Name:  "Коммунсервис",
		Email: "no-reply@yourdomain.com",
	})
	for _, e := range to {
		msg.AddRecipient(mailersend.Recipient{Name: "", Email: e})
	}
	msg.SetSubject(n.Title)
	msg.SetText(n.Message)
	msg.SetHTML("<p>" + n.Message + "</p>")
	if _, err := m.client.Email.Send(ctx, msg); err != nil {
		log.Printf("MailerSend error: %v", err)
	}
}
