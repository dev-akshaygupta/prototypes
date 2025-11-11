package main

import (
	"log"
	"net/http"

	"github.com/IBM/sarama"
)

func main() {
	http.HandleFunc("/addUser", handler.addUserHandler)
	http.HandleFunc("/updateUser", updateUserHandler)
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

// Create Producer
func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushUserToQueue(topic string, message []byte) error {
	broker := []string{"localhost:9092"}

	producer, err := ConnectProducer(broker)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Pushing to %s", topic)

	return nil
}
