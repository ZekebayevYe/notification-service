package tests

import (
	"context"
	"github.com/ZekebayevYe/notification-service/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockRepo struct {
	added   []string
	removed []string
}

func (m *mockRepo) AddSubscriber(ctx context.Context, email string) error {
	m.added = append(m.added, email)
	return nil
}
func (m *mockRepo) RemoveSubscriber(ctx context.Context, email string) error {
	m.removed = append(m.removed, email)
	return nil
}
func (m *mockRepo) ListSubscribers(ctx context.Context) ([]string, error) {
	return []string{"a@a.com"}, nil
}
func (m *mockRepo) SaveNotification(ctx context.Context, n app.Notification) error {
	return nil
}

type mockMailer struct {
	calls []app.Notification
}

func (m *mockMailer) SendNotification(n app.Notification, to []string) {
	m.calls = append(m.calls, n)
}

func TestCreateNotification(t *testing.T) {
	repo := &mockRepo{}
	mail := &mockMailer{}
	svc := app.NewService(repo, mail, 1*time.Minute)

	notification := app.Notification{
		Title:   "T",
		Message: "M",
		SendAt:  time.Now().Unix(),
	}

	err := svc.CreateNotification(context.Background(), notification)
	assert.NoError(t, err)

	mail.SendNotification(notification, []string{"a@a.com"})

	assert.Len(t, mail.calls, 1)
}
