package main

import (
    "fmt"
    "time"
    "context"
    "github.com/satori/go.uuid"
    opt "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/notification"

    "storage/instance/log"
    "storage/instance/minio"
)

func main() {
    log.Info("start", "Loading")

    client := minio.GetInstance()
    log.Debug("minio", "client:", fmt.Sprintf("%+v", client))

    ctx := context.Background()
    arndoc := notification.NewArn("minio", "sqs", "", "DOCUMENTS", "kafka")
    arnsig := notification.NewArn("minio", "sqs", "", "SIGNATURES", "kafka")

    for {
        date := time.Now() 
        bname := fmt.Sprintf("%d%02d%02d%02d%02d", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute())
        log.Debug("bucket", "name: ", bname)

        is_exists, err := client.BucketExists(ctx, bname)
        if err != nil {
            log.Error("error", "exists bucket: ", bname, ", error:", err)
            continue
        } else if !is_exists {
            err = client.MakeBucket(ctx, bname, opt.MakeBucketOptions{})
            if err != nil {
                log.Error("error", "create bucket:", bname, ", error:", err)
                continue
            }
            log.Debug("create", "bucket:", bname)

            queuedoc := notification.NewConfig(arndoc)
            queuedoc.AddEvents(notification.ObjectCreatedAll)
            queuedoc.AddFilterSuffix(".xml")

            queuesig := notification.NewConfig(arnsig)
            queuesig.AddEvents(notification.ObjectCreatedAll)
            queuesig.AddFilterSuffix(".sig")

            config := notification.Configuration{}
            config.AddQueue(queuesig)
            config.AddQueue(queuedoc)

            err = client.SetBucketNotification(ctx, bname, config)
            if err != nil {
                log.Error("error", "notification signature:", bname, ", error:", err)
                continue
            }
        }

        id := uuid.NewV4()
    
        err = minio.Upload("../document.xml", bname, fmt.Sprintf("%s.xml", id), "text/xml")
        if err != nil {
            log.Error("error", "upload", err)
            continue
        }

        time.Sleep(time.Second * 5)
    }

    log.Info("stop", "Loading")
}

