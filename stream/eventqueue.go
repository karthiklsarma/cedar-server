package stream

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/golang/protobuf/proto"
	"github.com/karthiklsarma/cedar-logging/logging"
	"github.com/karthiklsarma/cedar-schema/gen"
)

type IEventQueue interface {
	Connect() error
	EmitLocation(location *gen.Location) error
}

type EventQueue struct {
	hub                      *eventhub.Hub
	stream_connection_string string
}

func (q *EventQueue) GetHub() *eventhub.Hub {
	return q.hub
}

func (q *EventQueue) SetStreamConnectionString(constr string) {
	q.stream_connection_string = constr
}

func (eventQueue *EventQueue) Connect() error {
	if len(eventQueue.stream_connection_string) == 0 {
		eventQueue.stream_connection_string = os.Getenv(STREAM_CONN_ENV)
		if len(eventQueue.stream_connection_string) == 0 {
			err := "stream connection string is not set as environment variable."
			logging.Fatal(err)
			return fmt.Errorf(err)
		}
	}

	logging.Info("Initializing event hub...")
	var err error
	if eventQueue.hub, err = eventhub.NewHubFromConnectionString(eventQueue.stream_connection_string); err != nil {
		logging.Error(fmt.Sprintf("error initiating eventhub. error: %v", err))
		return err
	}

	return nil
}

func (eventQueue *EventQueue) EmitLocation(location *gen.Location) error {
	if location == nil {
		logging.Error("invalid input location")
		return errors.New("location is empty")
	}

	errorChan := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	go func(errorChan chan error) {
		logging.Info(fmt.Sprintf("received location: %v", location))
		location_bytes, err := proto.Marshal(location)
		if err != nil {
			errorChan <- err
			return
		}

		if err = eventQueue.hub.Send(ctx, eventhub.NewEvent(location_bytes)); err != nil {
			logging.Error(fmt.Sprintf("Something went wrong while sending msg to eventhub: %v", err))
			errorChan <- err
			return
		}
		logging.Info(fmt.Sprintf("successfully sent message %v to eventhub", location))
		close(errorChan)
	}(errorChan)

	return <-errorChan
}
