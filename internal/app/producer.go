package app

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
}

func NewKafkaProducer(brokers string, logger *logrus.Logger) (*Producer, error) {
	log.Printf("Инициализация Kafka producer. Используем брокеры: %s", brokers)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		log.Printf("Ошибка подключения к Kafka. Брокеры: %s, Ошибка: %v", brokers, err)
		return nil, err
	}

	log.Printf("Kafka producer успешно инициализирован. Подключен к брокерам: %s", brokers)
	return &Producer{producer: producer}, nil
}

func (p *Producer) SendMessage(topic string, message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		log.Printf("Ошибка закрытия Kafka producer: %v", err)
	} else {
		log.Println("Kafka producer закрыт")
	}
}
