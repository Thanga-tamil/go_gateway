package service

import (
	"fmt"
	"log"

	"github.com/Thanga-tamil/noway_service/internal/config"
)

var Chan = make(chan []byte)

func AmqpConsumer(msg chan []byte){
	
	msgs, err := config.AmqpChannel.Consume(
		"data", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		log.Printf("Received a message: %s\n", d.Body)
		go func(){ Chan <-d.Body }()
		go amqpConsumer()
	}
}

func amqpConsumer(){
	msg := <- Chan
	fmt.Printf("damn: %s\n", msg)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
