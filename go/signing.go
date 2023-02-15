package main

import (
    "fmt"
    "time"
    "regexp"
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"

    "storage/instance/log"
    "storage/instance/minio"
)

type Message struct {
    EventName string        `json:"EventName"`
    Key string              `json:"Key"`
}

func main() {
    log.Info("start", "Signing")

    ctx := context.Background()
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{"localhost:9092"},
        GroupID:   "signing",
        Topic:     "documents",
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

        msg := Message{}
        json.Unmarshal(event.Value, &msg)
        log.Info("event", "EventName:", msg.EventName, "Key:", msg.Key)

         re_bucket := regexp.MustCompile(`^\d+`)
         bucket := re_bucket.FindString(msg.Key)

        re_uuid := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
        uuid := re_uuid.FindString(msg.Key)

        log.Info("event", "bucket:", bucket, "uuid:", uuid)

        err = minio.Upload("../document.sig", bucket, fmt.Sprintf("%s.sig", uuid), "application/octet-stream")
        if err != nil {
            log.Error("error", "upload", err)
            continue
        }
    }

    log.Info("stop", "Signing")
}
