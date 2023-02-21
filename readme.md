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

### Delete topic

    $ docker exec -it kafka kafka-topics --bootstrap-server kafka:9092 --delete -topic signatures
    $ docker exec -it kafka kafka-topics --bootstrap-server kafka:9092 --delete -topic documents

## Subscribing to events

    $ docker exec -it kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic signatures
    $ docker exec -it kafka kafka-console-consumer --bootstrap-server kafka:9092 --topic documents

or

    $ kafkacat -C -b localhost:9092 -t signatures
    $ kafkacat -C -b localhost:9092 -t documents

# Uploading documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./loading.go

Creates a bucket every minute. The format of the bucket name is YYYYMMDDHHMM.
Uploads the {uuid}.xml file every 5 seconds

# Signing uploaded documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./signing.go

Downloads signature files. File format {uuid}.sig

# Archiving uploaded documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./archiving.go

# Deleting documents

    $ cd ./go
    $ export STORAGE_LOG_LEVEL=Debug; go run ./deleting.go

# References

 * [Multi-Cloud Object Storage](https://min.io/)
 * [MinIO Go Client API](https://min.io/docs/minio/linux/developers/go/API.html)
 * [kafka-go](https://github.com/segmentio/kafka-go)

