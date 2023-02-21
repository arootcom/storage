package main

import (
    "fmt"
    "context"
    opt "github.com/minio/minio-go/v7"

    "storage/instance/log"
    "storage/instance/minio"
    "storage/instance/minio/archive"
)

func main() {
    log.Info("start", "Deleting")

    delete_minio()
    delete_archive()

    log.Info("stop", "Deleting")
}

func delete_minio() {
    client := minio.GetInstance()

    ctx := context.Background()

    buckets, err := client.ListBuckets(ctx)
    if err != nil {
        log.Error("error", "list:", err)
        return
    }

    for _, bucket := range buckets {
        log.Debug("bucket", "bucket:", fmt.Sprintf("%+v", bucket))

        list := client.ListObjects(ctx, bucket.Name, opt.ListObjectsOptions{})

        for obj := range list {
            err =client.RemoveObject(ctx, bucket.Name, obj.Key, opt.RemoveObjectOptions{
                GovernanceBypass: true,
            })
            if err != nil {
                log.Error("error", "delete:", err)
                continue
            }
            log.Debug("deleted", "object:", obj.Key)
        }

        err = client.RemoveBucket(ctx, bucket.Name)
        if err != nil {
            log.Error("error", "delete:", err)
            continue
        }
        log.Info("deleted", "bucket:", bucket.Name)
    }
}

func delete_archive() {
    client := archive.GetInstance()

    ctx := context.Background()

    buckets, err := client.ListBuckets(ctx)
    if err != nil {
        log.Error("error", "list:", err)
        return
    }

    for _, bucket := range buckets {
        log.Debug("bucket", "bucket:", fmt.Sprintf("%+v", bucket))

        list := client.ListObjects(ctx, bucket.Name, opt.ListObjectsOptions{
            Recursive: true,
            WithVersions: true,
        })

        for obj := range list {
            err =client.RemoveObject(ctx, bucket.Name, obj.Key, opt.RemoveObjectOptions{
                ForceDelete: true,
                GovernanceBypass: true,
            })
            if err != nil {
                log.Error("error", "delete:", err)
                continue
            }
            log.Debug("deleted", "object:", obj.Key)
        }

        err = client.RemoveBucket(ctx, bucket.Name)
        if err != nil {
            log.Error("error", "delete:", err)
            continue
        }
        log.Info("deleted", "bucket:", bucket.Name)
    }
}
