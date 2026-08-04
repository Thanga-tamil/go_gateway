package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var AmqpChannel *amqp.Channel

func InitRabbitMQConnection() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	AmqpChannel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
