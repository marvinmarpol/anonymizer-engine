package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marvinmarpol/golang-boilerplate/internal/protection/entity"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/service"
	"github.com/sirupsen/logrus"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Server struct {
	service service.Services
}

func NewServer(service service.Services) *Server {
	return &Server{service}
}

func (s *Server) RotateKeys(msg *message.Message) error {
	// Acknowledge the message
	msg.Ack()

	start := time.Now() // Capture start time

	// Define var to store the unmarshalled JSON payload
	var payload entity.RotatePayload

	// Unmarshal the JSON payload into a map
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		logrus.Error("Failed to unmarshal JSON payload")
		return err
	}

	result, err := s.service.RotateKeys(context.Background(), payload)
	if err != nil {
		logrus.WithField("err", err).Error("Failed to rotate keys")
		return err
	}

	elapsed := time.Since(start) // Calculate elapsed time

	logrus.Info(fmt.Sprintf("API took %s to process", elapsed))
	logrus.Info("Total rotated key: ", result)
	return nil
}
