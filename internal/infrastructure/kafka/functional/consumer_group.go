package functional

import (
	"Avito/internal/storage/repository"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"strconv"
	"strings"
)

type paymentMessage struct {
	AnswerURL string
	Method    string
	Success   bool
}

type ConsumerGroup struct {
	ready chan bool
	Repo  repository.BannerRepo
}

func NewConsumerGroup(implementation repository.BannerRepo) ConsumerGroup {
	return ConsumerGroup{
		ready: make(chan bool),
		Repo:  implementation,
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			pm := paymentMessage{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				fmt.Println("Consumer group error", err)
			}

			parts := strings.Split(pm.AnswerURL, "/")
			// Возвращаем первые два компонента пути (или всё, что есть, если менее двух).
			url := parts[2]
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				panic(err)
			}

			if url == "feature" {
				consumer.Repo.DeleteByFeatureIDHandler(context.Background(), int64(id))
			} else {
				consumer.Repo.DeleteByTagIDHandler(context.Background(), int64(id))
			}

			log.Printf("Message claimed: URL = %v,Message claimed: Method = %v, timestamp = %v, topic = %s",
				pm.AnswerURL,
				pm.Method,
				message.Timestamp,
				message.Topic,
			)

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
