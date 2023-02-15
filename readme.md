# Storage

# Start

    $ cd ./manifest
    $ docker-compose up -d

# MINIO

    Open in a browser 
    http://localhost:9001/
    Using credentials: minioadmin/minioadmin.

# Kafka

## List of topics

    $ docker exec -it kafka kafka-topics --list --bootstrap-server kafka:9092

## Subscribing to events

    $ docker exec -it kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic notifications

or 
    
    $ kafkacat -C -b localhost:9092 -t notifications

# Uploading documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./loading.go

Creates a bucket every minute. The format of the bucket name is YYYYMMDDHHMM.
Downloads two files with xml and sig extensions every 5 seconds. The file name in uuid format.

# Deleting documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./deleting.go 


# References

 * [Multi-Cloud Object Storage](https://min.io/)
 * [MinIO Go Client API](https://min.io/docs/minio/linux/developers/go/API.html)
 * [kafka-go](https://github.com/segmentio/kafka-go)

