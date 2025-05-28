package messaging

import (
	"encoding/json"
	"github.com/ZekebayevYe/notification-service/internal/model"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

var nc *nats.Conn

func Init() {
	url := os.Getenv("NATS_URL")
	var err error
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatalf("nats.Connect: %v", err)
	}
}

func PublishNotification(n model.Notification) {
	data, _ := json.Marshal(n)
	nc.Publish("notifications", data)
}

func SubscribeNotifications(handler func(model.Notification)) {
	nc.QueueSubscribe("notifications", "workers", func(msg *nats.Msg) {
		var n model.Notification
		if err := json.Unmarshal(msg.Data, &n); err != nil {
			log.Printf("unmarshal: %v", err)
			return
		}
		handler(n)
	})
}
