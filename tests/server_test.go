package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/google/uuid"
	"github.com/karthiklsarma/cedar-schema/gen"
	"github.com/karthiklsarma/cedar-server/stream"
)

func TestEventQueueAdd(t *testing.T) {
	eventQ := &stream.EventQueue{}
	eventQ.SetStreamConnectionString(test_stream_connection_string)
	eventQ.Connect()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	hub := eventQ.GetHub()
	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		fmt.Println("Failed to obtain runtime information for eventhub.")
		t.Fail()
	}

	received := false
	for _, partition := range runtimeInfo.PartitionIDs {
		hub.Receive(ctx, partition, func(ctx context.Context, event *eventhub.Event) error {
			fmt.Println("Received Location...")
			received = true
			return nil
		}, eventhub.ReceiveWithLatestOffset())
	}

	test_id := uuid.New().String()
	location := &gen.Location{Id: test_id, Lat: 1.1, Lng: 1.1, Timestamp: uint64(time.Now().Unix()), Device: "test_device"}
	eventQ.EmitLocation(location)
	time.Sleep(2 * time.Second)
	if !received {
		t.Fail()
	}
}
