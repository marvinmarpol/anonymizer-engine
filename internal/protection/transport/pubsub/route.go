package pubsub

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type MiddlewareFunc message.HandlerMiddleware

type Route interface {
	RotateKeys(msg *message.Message) error
}

func RegisterRoutes(r *message.Router, si Route, subsriber message.Subscriber, Middlewares []MiddlewareFunc) {
	// register middlewares
	for _, midleware := range Middlewares {
		r.AddMiddleware(message.HandlerMiddleware(midleware))
	}

	// just for debug, we are printing all messages received on `incoming_messages_topic`
	r.AddNoPublisherHandler(
		"rotateKeyHandler",
		"rotateKeyTopic",
		subsriber,
		si.RotateKeys,
	)
}
