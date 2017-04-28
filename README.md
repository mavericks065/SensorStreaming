This project's goal is to publish fake sensor data to a messaging system: either RabbitMQ, Kafka or Flink.

First to have it working you need to run :

go get -u github.com/jessevdk/go-flags
go get github.com/Shopify/sarama
go get -u "github.com/eapache/queue"
go get -u "github.com/rcrowley/go-metrics"
go get -u "github.com/klauspost/crc32"
go get -u "github.com/pierrec/lz4"
go get -u "github.com/golang/snappy"

go build

go run /sensor/Sensor.go ... // with the many options
