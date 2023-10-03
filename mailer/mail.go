package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"rehab/internal/pkg/models"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func receive(channel string, template string) *redis.Message {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),      // Replace with your Redis server address
		Password: os.Getenv("REDIS_PASSWORD"), // Set if you have a password for your Redis server
		DB:       0,                           // Select the appropriate Redis database
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)

	}

	var data models.Email
	pubsub := client.Subscribe(context.Background(), channel)
	_, err = pubsub.Receive(context.Background())
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println("Message Received from channel " + channel)
		err = json.Unmarshal([]byte(msg.Payload), &data)
		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			continue
		}
		if err = SendEmail(data.MailAddress, &data, template); err != nil {
			fmt.Println(err)
		}

	}
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}

func SendEmail(mail_address string, data *models.Email, templateName string) error {

	from := os.Getenv("EMAIL_FROM")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpUser := os.Getenv("SMTP_USER")
	to := mail_address
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return err
	}

	template.ExecuteTemplate(&body, templateName, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	int_smtpPort, _ := strconv.Atoi(smtpPort)
	d := gomail.NewDialer(smtpHost, int_smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return err
}

func main() {
	// fmt.Println("asd")
	// go receive("Accepted_Bids", "Validateorder.html")
	go receive("registration", "verificationCode.html")
	// go receive("User_reset_password", "resetPassword.html")
	select {}
}
