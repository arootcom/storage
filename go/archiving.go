package main

import (
    "fmt"
    "time"
    "regexp"
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
    opt "github.com/minio/minio-go/v7"

    "storage/instance/log"
    "storage/instance/minio"
    "storage/instance/minio/archive"
)

type Message struct {
    EventName string        `json:"EventName"`
    Key string              `json:"Key"`
}

func main() {
    log.Info("start", "Archiving")

    client := minio.GetInstance()
    arc := archive.GetInstance()

    ctx := context.Background()
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:   []string{"localhost:9092"},
        GroupID:   "archiving",
        Topic:     "signatures",
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

        is_exists, err := arc.BucketExists(ctx, bucket)
        if err != nil {
            log.Error("error", "exists bucket: ", bucket, ", error:", err)
            continue
        } else if !is_exists {
            err = arc.MakeBucket(ctx, bucket, opt.MakeBucketOptions{})
            if err != nil {
                log.Error("error", "create bucket:", bucket, ", error:", err)
                continue
            }

            err = arc.EnableVersioning(ctx, bucket)
            if err != nil {
                log.Error("error", "versioning bucket:", bucket, ", error:", err)
                continue
            }

            log.Debug("create", "bucket:", bucket)
        }

        for _, ext := range []string{"sig", "xml"} {
            file_name := fmt.Sprintf("%s.%s", uuid, ext)

            stat, err := client.StatObject(ctx, bucket, file_name, opt.StatObjectOptions{})
            if err != nil {
                log.Error("error", "stat:", err)
                break
            }
            log.Debug("archiving", "stat:", fmt.Sprintf("%+v", stat))

            tags, err := client.GetObjectTagging(ctx, bucket, file_name, opt.GetObjectTaggingOptions{})
            if err != nil {
                log.Error("error", "tags:", err)
                break
            }
            log.Debug("archiving", "tags:", fmt.Sprintf("%+v", tags))

            reader, err := client.GetObject(ctx, bucket, file_name, opt.GetObjectOptions{})
            if err != nil {
                log.Error("error", "reader:", err)
                break
            }
            defer reader.Close()

            object, err := arc.PutObject(ctx, bucket, file_name, reader, stat.Size,
                opt.PutObjectOptions{
                    ContentType: stat.ContentType,
                    UserMetadata: stat.UserMetadata,
                },
            )
            if err != nil {
                log.Error("error", "create:", err)
                break
            }
            log.Debug("archiving", "object:", fmt.Sprintf("%+v", object))

            err = arc.PutObjectTagging(ctx, bucket, file_name, tags, opt.PutObjectTaggingOptions{})
            if err != nil {
                log.Error("error", "tags", err)
                continue
            }

            log.Info("archiving",
                "file:", fmt.Sprintf("/%s/%s", bucket, file_name),
                ", size:", stat.Size, ", content-type:", stat.ContentType, ", meta:", stat.UserMetadata,
            )

            if ext == "xml" {
                dt := time.Now()
                archived := fmt.Sprintf("%04d%02d%02dT%02d%02d%02d",                         //  ISO 8601
                    dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second(),
                )
                tags.Set("Archived", archived)

                log.Debug("archiving", "tags:", fmt.Sprintf("%+v", tags))

                err = client.PutObjectTagging(ctx, bucket, file_name, tags, opt.PutObjectTaggingOptions{})
                if err != nil {
                    log.Error("error", "tags", err)
                    continue
                }
            }
        }

    }

    log.Info("stop", "Archiving")
}

