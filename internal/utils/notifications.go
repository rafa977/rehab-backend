package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	//"github.com/gin-gonic/gin"
	//"github.com/go-redis/redis"

	"github.com/go-redis/redis/v8"
	"github.com/rehab-backend/internal/pkg/models"
)

type RedisHandler interface {
	Redis_Mailer(string, string, string)
	Redis_Mailer_ui(string, string, string)
	// Notifications(uint, models.Notification, string)
	Add_mail_to_Redis_queue(string, string, string, string)
}

type redisRepository struct {
	connection *redis.Client
}

func NewredisRepository() RedisHandler {
	return &redisRepository{
		connection: ConnectonRedis(),
	}
}

func ConnectonRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_url"),      // Replace with your Redis server address
		Password: os.Getenv("redis_password"), // Set if you have a password for your Redis server
		DB:       0,                           // Select the appropriate Redis database
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return nil
	}
	return client
}

func (client *redisRepository) Add_mail_to_Redis_queue(mail string, url string, subject string, channel string) {
	emailData := models.Email{
		MailAddress: mail,
		URL:         url,
		Subject:     subject,
	}
	payload, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println("Failed to marshal struct to JSON:", err)
		return
	}

	err = client.connection.Publish(context.Background(), channel, string(payload)).Err()
	if err != nil {
		fmt.Println("Failed to publish message on Redis:", err)
		return
	}
	log.Println("Message send on channel " + channel)
}

func (client *redisRepository) Redis_Mailer(mail string, token string, version string) {
	emailData := models.Email{
		MailAddress: mail,
		URL:         os.Getenv("CLIENT_ORIGIN") + os.Getenv(version+"_URL"),
		TOKEN:       token,
		Subject:     os.Getenv(version + "_Subject"),
	}
	payload, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println("Failed to marshal struct to JSON:", err)
		return
	}

	err = client.connection.Publish(context.Background(), os.Getenv(version+"_Channel"), string(payload)).Err()
	if err != nil {
		fmt.Println("Failed to publish message on Redis:", err)
		return
	}
	log.Println("Message send on channel " + version)
}

func (client *redisRepository) Redis_Mailer_ui(mail string, token string, version string) {
	emailData := models.Email{
		MailAddress: mail,
		URL:         os.Getenv("SOURCE") + os.Getenv(version+"_URL"),
		TOKEN:       token,
		Subject:     os.Getenv(version + "_Subject"),
	}
	payload, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println("Failed to marshal struct to JSON:", err)
		return
	}

	err = client.connection.Publish(context.Background(), os.Getenv(version+"_Channel"), string(payload)).Err()
	if err != nil {
		fmt.Println("Failed to publish message on Redis:", err)
		return
	}
	log.Println("Message send on channel " + version)
}

// func (client *redisRepository) Notifications(id uint, notification models.Notification, version string) {
// 	var err error
// 	if err != nil {
// 		fmt.Println("Failed to marshal struct to JSON:", err)
// 		return
// 	}
// 	notification.ProfileID = id
// 	if _, err = customer.NewUserRepository().AddNotification(notification.Subject); err != nil {
// 		log.Default()
// 	}
// 	err = client.connection.Publish(context.Background(), os.Getenv(version+"_Channel"), notification.Subject).Err()
// 	if err != nil {
// 		fmt.Println("Failed to publish message on Redis:", err)
// 		return
// 	}
// 	fmt.Println("message published")
// }
