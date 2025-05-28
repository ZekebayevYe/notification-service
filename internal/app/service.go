package app

import (
	"context"
	"errors"
	"github.com/ZekebayevYe/notification-service/internal/model"
	"time"
)

var (
	ErrAlreadySubscribed = errors.New("уже подписан")
)

type Notification = model.Notification

type Service struct {
	repo     Repository
	cacheTTL time.Duration
	mailer   Mailer
}

type Repository interface {
	AddSubscriber(context.Context, string) error
	RemoveSubscriber(context.Context, string) error
	ListSubscribers(context.Context) ([]string, error)
	SaveNotification(context.Context, Notification) error
}

type Mailer interface {
	SendNotification(Notification, []string)
}

func NewService(r Repository, m Mailer, cacheTTL time.Duration) *Service {
	return &Service{repo: r, mailer: m, cacheTTL: cacheTTL}
}

func (s *Service) Subscribe(ctx context.Context, email string) error {
	return s.repo.AddSubscriber(ctx, email)
}

func (s *Service) Unsubscribe(ctx context.Context, email string) error {
	return s.repo.RemoveSubscriber(ctx, email)
}

func (s *Service) CreateNotification(ctx context.Context, n Notification) error {
	if err := s.repo.SaveNotification(ctx, n); err != nil {
		return err
	}
	return nil
}
