package messaging

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/api/option"
)

type GooglePubSub struct {
	publisher  message.Publisher
	Subscriber message.Subscriber
	Router     message.Router
}

func DefaultLogger() watermill.LoggerAdapter {
	return watermill.NewStdLogger(false, false)
}

func NewGooglePubSub(projectID, filepath string, logger watermill.LoggerAdapter) (*GooglePubSub, error) {
	credOption := []option.ClientOption{option.WithCredentialsFile(filepath)}
	pubConfig := googlecloud.PublisherConfig{
		ProjectID:     projectID,
		ClientOptions: credOption,
	}

	publisher, err := googlecloud.NewPublisher(pubConfig, logger)
	if err != nil {
		return nil, err
	}

	subConfig := googlecloud.SubscriberConfig{
		ProjectID:     projectID,
		ClientOptions: credOption,
	}

	subscriber, err := googlecloud.NewSubscriber(subConfig, logger)
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(message.RouterConfig{
		CloseTimeout: time.Second * 60,
	}, logger)
	if err != nil {
		return nil, err
	}

	return &GooglePubSub{
		publisher:  publisher,
		Subscriber: subscriber,
		Router:     *router,
	}, nil
}

func (g *GooglePubSub) Publish(topic string, msg *message.Message) error {
	return g.publisher.Publish(topic, msg)
}

func (g *GooglePubSub) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return g.Subscriber.Subscribe(ctx, topic)
}

func (g *GooglePubSub) Close() error {
	if err := g.publisher.Close(); err != nil {
		return err
	}
	return g.Subscriber.Close()
}
