package minio

import (
    "os"
    "fmt"
    "sync"
    "context"
    "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

    "storage/instance/log"
)

var once sync.Once
var instance *minio.Client = nil

//
func GetInstance() *minio.Client {
    once.Do(func() {
        var err error
        instance, err = minio.New("localhost:9000", &minio.Options{
            Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
            Secure: false,
        })
        if err != nil {
            panic(err)
        }
    })
    return instance
}

//
func Upload (filename string, bucket string, key string, contentType string, meta map[string]string) error {
    reader, err := os.Open(filename)
    defer reader.Close()
    if err != nil {
        return err
    }

    stat, err := reader.Stat()
    if err != nil {
        return err
    }

    client := GetInstance()
    ctx := context.Background()

    object, err := client.PutObject(ctx, bucket, key, reader, stat.Size(),
        minio.PutObjectOptions{
            ContentType: contentType,
            UserMetadata: meta,
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
