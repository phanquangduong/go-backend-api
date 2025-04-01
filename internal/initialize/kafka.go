package initialize

import (
	"go/go-backend-api/global"
	"log"

	"github.com/segmentio/kafka-go"
)

var KafkaProducer *kafka.Writer

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:19094"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}

}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %v", err)
	}
}
