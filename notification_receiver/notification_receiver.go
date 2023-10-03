package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

	go receive("hells", ctx, client)
	select {}
}

func receive(channel string, ctx context.Context, client *redis.Client) *redis.Message {
	pubsub := client.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)

	}

	// defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(msg.Payload)

	}
}
