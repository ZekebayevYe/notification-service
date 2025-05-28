package repository

import (
	"context"
	"github.com/ZekebayevYe/notification-service/internal/app"
	"github.com/ZekebayevYe/notification-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Subscriber struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Email string             `bson:"email"`
}

type NotificationEntity struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title"`
	Message string             `bson:"message"`
	SendAt  int64              `bson:"send_at"`
}

type NotificationRepo struct {
	db *mongo.Database
}

func NewNotificationRepo(client *mongo.Client) *NotificationRepo {
	return &NotificationRepo{db: client.Database("notifications")}
}

func (r *NotificationRepo) AddSubscriber(ctx context.Context, email string) error {
	coll := r.db.Collection("subscribers")
	_, err := coll.InsertOne(ctx, Subscriber{Email: email})
	if mongo.IsDuplicateKeyError(err) {
		return app.ErrAlreadySubscribed
	}
	return err
}

func (r *NotificationRepo) RemoveSubscriber(ctx context.Context, email string) error {
	coll := r.db.Collection("subscribers")
	_, err := coll.DeleteOne(ctx, bson.M{"email": email})
	return err
}

func (r *NotificationRepo) ListSubscribers(ctx context.Context) ([]string, error) {
	coll := r.db.Collection("subscribers")
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var subs []string
	for cur.Next(ctx) {
		var s Subscriber
		_ = cur.Decode(&s)
		subs = append(subs, s.Email)
	}
	return subs, nil
}

func (r *NotificationRepo) SaveNotification(ctx context.Context, n model.Notification) error {
	coll := r.db.Collection("notifications")
	_, err := coll.InsertOne(ctx, NotificationEntity{
		Title:   n.Title,
		Message: n.Message,
		SendAt:  n.SendAt,
	})
	return err
}
