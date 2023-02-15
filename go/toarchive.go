package main

import (
    "fmt"
    "time"
    "context"
    "github.com/segmentio/kafka-go"

    "storage/instance/log"
)

func main() {
    log.Info("start", "Listening")

    ctx := context.Background()
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{"localhost:9092"},
        GroupID:   "toarchive",
        Topic:     "notifications",
        MinBytes:  10e3, // 10KB
        MaxBytes:  10e6, // 10MB
    })

    for {
        event, err := reader.ReadMessage(ctx)
        if err != nil {
            time.Sleep(time.Second * 1)
        }

        log.Info("event", 
            fmt.Sprintf(
                "message at topic/partition/offset %v/%v/%v: %s = %s\n", 
                event.Topic, event.Partition, event.Offset, string(event.Key), string(event.Value),
            ),
        )
    }

    log.Info("stop", "Listening")
}

