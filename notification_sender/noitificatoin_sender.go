package main

import (
	"context"
	"fmt"
	"rehab/internal/pkg/models"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("Go Redis Tutorial")

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

	var not models.Notification
	not.Subject = "This is a not"
	notifications(not, ctx, client)

	// fmt.Println(val)

}

func notifications(notification models.Notification, ctx context.Context, client *redis.Client) {

	err := client.Publish(ctx, "registration", notification.Subject).Err()
	if err != nil {
		fmt.Println("Failed to publish message on Redis:", err)
		return
	}
	fmt.Println("message published")
}
