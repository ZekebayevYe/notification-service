package main

import (
	"context"
	"github.com/ZekebayevYe/notification-service/internal/app"
	"github.com/ZekebayevYe/notification-service/internal/cache"
	"github.com/ZekebayevYe/notification-service/internal/db"
	"github.com/ZekebayevYe/notification-service/internal/email"
	pb "github.com/ZekebayevYe/notification-service/internal/grpc"
	"github.com/ZekebayevYe/notification-service/internal/handler"
	"github.com/ZekebayevYe/notification-service/internal/messaging"
	"github.com/ZekebayevYe/notification-service/internal/model"
	"github.com/ZekebayevYe/notification-service/internal/repository"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	client := db.Connect()
	messaging.Init()
	ttl, _ := time.ParseDuration(os.Getenv("CACHE_TTL"))
	cache.Init(ttl)
	repo := repository.NewNotificationRepo(client)
	mailer := email.NewMailer()
	svc := app.NewService(repo, mailer, ttl)
	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, handler.NewGRPCServer(svc))
	log.Printf("gRPC server listening on %s …", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Serve error: %v", err)
	}

	messaging.SubscribeNotifications(func(n model.Notification) {
		subs, err := repo.ListSubscribers(context.Background())
		if err != nil {
			return
		}
		mailer.SendNotification(n, subs)
	})

}
