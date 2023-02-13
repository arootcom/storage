package main

import (
    "os"
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

            arn := notification.NewArn("minio", "sqs", "", "_", "kafka")
            queue := notification.NewConfig(arn)
            queue.AddEvents(notification.ObjectCreatedAll)
            queue.AddFilterSuffix(".sig")

            config := notification.Configuration{}
            config.AddQueue(queue)

            err = client.SetBucketNotification(ctx, bname, config)
            if err != nil {
                log.Error("error", "notification:", bname, ", error:", err)
                continue
            }
        }

        id := uuid.NewV4()
    
        err = upload("../document.xml", bname, fmt.Sprintf("%s.xml", id), "text/xml")
        if err != nil {
            log.Error("error", "upload", err)
            continue
        }
        
        err = upload("../document.sig", bname, fmt.Sprintf("%s.sig", id), "application/octet-stream")
        if err != nil {
            log.Error("error", "upload", err)
            continue
        }

        time.Sleep(time.Second * 5)
    }

    log.Info("stop", "Loading")
}

func upload (filename string, bucket string, key string, contentType string) error {
    reader, err := os.Open(filename)
    defer reader.Close()
    if err != nil {
        return err
    }
    log.Debug("upload", "open", filename)

    stat, err := reader.Stat()
    if err != nil {
        return err
    }
    log.Debug("upload", "size", stat.Size())

    client := minio.GetInstance()
    ctx := context.Background()

    object, err := client.PutObject(ctx, bucket, key, reader, stat.Size(),
        opt.PutObjectOptions{
            ContentType: contentType,
        },
    )
    if err != nil {
        return err
    }
    log.Debug("upload", "object:", object)
    
    log.Debug("upload", 
        "from", filename, "size", stat.Size(), 
        "to", fmt.Sprintf("/%s/%s", bucket, key), "content/type", contentType,
    )
    return nil
}
