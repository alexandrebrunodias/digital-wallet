package kafka

import (
	"encoding/json"
	"errors"
	"github.com/alexandrebrunodias/wallet-core/pkg/events"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strings"
	"sync"
)

type Producer struct {
	ConfigMap    *ckafka.ConfigMap
	Topic        *string
	PartitionKey []byte
}

func NewKafkaProducer(configMap *ckafka.ConfigMap, topic string, partitionKey []byte) *Producer {
	if configMap == nil {
		panic(errors.New("configMap must not be null"))
	}
	if strings.TrimSpace(topic) == "" {
		panic(errors.New("topic must not be null"))
	}
	return &Producer{
		ConfigMap:    configMap,
		Topic:        &topic,
		PartitionKey: partitionKey,
	}
}

func (p *Producer) Send(event events.Event, wg *sync.WaitGroup) error {
	defer wg.Done()
	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err != nil {
		return err
	}

	json, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: p.Topic, Partition: ckafka.PartitionAny},
		Value:          json,
		Key:            p.PartitionKey,
	}
	err = producer.Produce(message, nil)
	if err != nil {
		return err
	}
	return nil
}
