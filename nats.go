package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsResponse func(string) string

func getNatsURL() string {
	if len(os.Getenv("NATS_URL")) > 0 {
		return os.Getenv("NATS_URL")
	}
	return nats.DefaultURL
}

func NatsSubscribe() {
	if os.Getenv("NATS_SUBSCRIBE") == "true" {
		ctx, _ := context.WithCancel(context.Background())
		go natsSubscribe(ctx, "silly-demo.hello", func(message string) string {
			return "I'm the silliest demo you ever saw. Nice to meet you."
		})
		go natsSubscribe(ctx, "ci.silly-demo", func(message string) string {
			return "Thanks for doing the CI/CD for me."
		})
		go natsSubscribe(ctx, "fibonacci.request", func(message string) string {
			number, err := strconv.Atoi(message)
			if err != nil {
				return fmt.Sprintf("%s is not a number", message)
			}
			return strconv.Itoa(calculateFibonacci(number))
		})
	}
}

func NatsPublishLoop() {
	if os.Getenv("NATS_PUBLISH") == "true" {
		for {
			natsPublish("ping", "Silly demo is here. Is there anyone else around? Say hi to on the silly-demo.hello channel.")
			time.Sleep(10 * time.Second)
		}
	}
}

func natsPublish(channel, message string) error {
	nc, err := nats.Connect(getNatsURL())
	if err != nil {
		return err
	}
	defer nc.Close()
	log.Printf("publishing message: %s\n", message)
	err = nc.Publish(channel, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func natsSubscribe(ctx context.Context, channel string, fn NatsResponse) {
	nc, err := nats.Connect(getNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	messages := make(chan *nats.Msg, 1000)
	subscription, err := nc.ChanSubscribe(channel, messages)
	if err != nil {
		log.Fatal("Failed to subscribe to subject:", err)
	}
	defer func() {
		subscription.Unsubscribe()
		close(messages)
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("exiting from the message subscriber")
			return
		case message := <-messages:
			log.Printf("received message: %s\n", string(message.Data))
			response := fn(string(message.Data))
			message.Respond([]byte(response))
		}
	}
}
