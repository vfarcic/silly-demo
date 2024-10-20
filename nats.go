package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type CloudEventData struct {
	Message string `json:"message"`
}

type CloudEventResponse func(context.Context, cloudevents.Event) error

// FIXME: Remove
// type NatsResponse func(string) string
type NatsResponse func(ctx context.Context, event cloudevents.Event) error

var natsURL string

func init() {
	if len(os.Getenv("NATS_URL")) > 0 {
		natsURL = os.Getenv("NATS_URL")
	}
	natsURL = nats.DefaultURL
}

func NatsSubscribe() {
	if os.Getenv("NATS_SUBSCRIBE") == "true" {
		ctx, _ := context.WithCancel(context.Background())
		go natsSubscribe(ctx, "silly-demo.hello", func(ctx context.Context, event cloudevents.Event) error {
			data := &CloudEventData{}
			if err := event.DataAs(data); err != nil {
				slog.Error(err.Error())
				return err
			}
			slog.Debug("message received", "message", data.Message)
			return nil
			// return "I'm the silliest demo you ever saw. Nice to meet you."
		})
		// go natsSubscribe(ctx, "ci.silly-demo", func(message string) string {
		// 	return "Thanks for doing the CI/CD for me."
		// })
		// go natsSubscribe(ctx, "fibonacci.request", func(message string) string {
		// 	number, err := strconv.Atoi(message)
		// 	if err != nil {
		// 		return fmt.Sprintf("%s is not a number", message)
		// 	}
		// 	return strconv.Itoa(CalculateFibonacci(number))
		// })
	}
}

func NatsPublish(eventType, subject, message string) error {
	sender, err := cenats.NewSender(natsURL, subject, cenats.NatsOptions())
	if err != nil {
		return err
	}
	defer sender.Close(context.Background())
	client, err := cloudevents.NewClient(sender)
	if err != nil {
		return err
	}
	e := getEvent(subject, message, eventType)
	if result := client.Send(context.TODO(), *e); cloudevents.IsUndelivered(result) {
		return fmt.Errorf("failed to send: %v", result)
	} else {
		slog.Debug("message sent", "message", message)
	}
	return nil
}

func NatsPublishLoop() {
	if os.Getenv("NATS_PUBLISH") == "true" {
		for {
			err := NatsPublish(
				"sent",
				"ping",
				"Silly demo is here. Is there anyone else around? Say hi on the silly-demo.hello channel.",
			)
			if err != nil {
				slog.Error(err.Error())
			}
			time.Sleep(1 * time.Minute)
		}
	}
}

func getEvent(subject, message, eventType string) *cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetType(fmt.Sprintf("live.devopstoolkit.silly-demo.%s.%s", subject, eventType))
	e.SetTime(time.Now())
	e.SetSource(fmt.Sprintf("https://devopstoolkit.live/demo/silly-demo/%s", subject))
	_ = e.SetData("application/json", &CloudEventData{
		Message: message,
	})
	return &e
}

func natsSubscribe(ctx context.Context, subject string, fn NatsResponse) {
	consumer, err := cenats.NewConsumer(natsURL, subject, cenats.NatsOptions())
	if err != nil {
		slog.Error(err.Error())
	}
	defer consumer.Close(context.Background())
	client, err := cloudevents.NewClient(consumer)
	if err != nil {
		slog.Error(err.Error())
	}
	for {
		if err := client.StartReceiver(ctx, fn); err != nil {
			slog.Error(err.Error())
		}
	}
	// if err != nil {
	// 	return err
	// }
	// e := getEvent(subject, message, eventType)
	// if result := client.Send(context.Background(), *e); cloudevents.IsUndelivered(result) {
	// 	return fmt.Errorf("failed to send: %v", result)
	// } else {
	// 	log.Printf("message accepted: %s", message)
	// }
	// return nil

	// FIXME: Remove
	// nc, err := nats.Connect(natsURL)
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// defer nc.Close()

	// messages := make(chan *nats.Msg, 1000)
	// subscription, err := nc.ChanSubscribe(subject, messages)
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// defer func() {
	// 	subscription.Unsubscribe()
	// 	close(messages)
	// }()
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("exiting from the message subscriber")
	// 		return
	// 	case message := <-messages:
	// 		log.Printf("received message: %s\n", string(message.Data))
	// 		response := fn(string(message.Data))
	// 		message.Respond([]byte(response))
	// 	}
	// }
}
