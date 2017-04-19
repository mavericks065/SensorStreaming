This project's goal is to publish fake sensor data to a messaging system: either RabbitMQ, Kafka or Flink.

First to have it working you need to run :

go get -u github.com/jessevdk/go-flags
go build
go run /sensor/Sensor.go ... // with the many options
