package minio

import (
    "sync"
    "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var once sync.Once
var instance *minio.Client = nil

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

