package abango

import (
	"log"
	"sync"

	"github.com/Shopify/sarama"
	e "github.com/dabory/abango-rest/etc"
)

var (
	KAFKA_CONN  string
	KAFKA_TOPIC string
)

func init() {
	KAFKA_CONN = XConfig["KafkaConnString"]
	KAFKA_TOPIC = XConfig["KafkaTopic"]
	KAFKA_TIMEOUT = XConfig["KafkaTimeout"]
}

func KafkaProducer(key string, message string, conCurr string) (int32, int64, error) {

	kfcf := sarama.NewConfig()
	kfcf.Producer.Retry.Max = 5
	kfcf.Producer.RequiredAcks = sarama.WaitForAll
	kfcf.Producer.Return.Successes = true

	if conCurr == "async" {
		if prd, err := sarama.NewAsyncProducer([]string{KAFKA_CONN}, kfcf); err == nil {
			prd.Input() <- &sarama.ProducerMessage{
				Topic: KAFKA_TOPIC,
				Key:   sarama.StringEncoder(key),     //[]byte doesn't work.
				Value: sarama.StringEncoder(message), //[]byte doesn't work.
			}
			return 0, 0, nil
		} else {
			return 0, 0, e.MyErr("QEJHDRTTRRW-Kafka-NewSyncProducer-End", err, true)
		}
	} else if conCurr == "sync" {
		if prd, err := sarama.NewSyncProducer([]string{conn}, kfcf); err == nil {
			msg := &sarama.ProducerMessage{
				Topic: KAFKA_TOPIC,
				Key:   sarama.StringEncoder(key),     //[]byte doesn't work.
				Value: sarama.StringEncoder(message), //[]byte doesn't work.
			}
			if part, offset, err := prd.SendMessage(msg); err == nil {
				return part, offset, nil
			} else {
				return 0, 0, e.MyErr("QEJIOPRTRRTRRW-Kafka-Sync-SendMessage", err, true)
			}
		} else {
			return 0, 0, e.MyErr("QEJHGTRSDRTTRRW-Kafka-NewSyncProducer-End", err, true)
		}
	} else {
		return 0, 0, e.MyErr("QEJHGTRSW-Kafka-ApiMethod Not available-End", nil, true)
	}
}

func KafkaConsumer() {

	// Create a new configuration for the consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify the list of brokers
	brokers := []string{KAFKA_CONN}

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalf("Error closing consumer: %s", err)
		}
	}()

	// Create a new consumer for topics
	// topics := KAFKA_TOPIC
	partitions, err := consumer.Partitions(KAFKA_TOPIC)
	if err != nil {
		log.Fatalf("Failed to get partitions: %s", err)
	}
	//lsjdljsdlfjas
	var wg sync.WaitGroup
	wg.Add(len(partitions))

	// Consume messages from each partition asynchronously
	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()

			// Create a new partition consumer
			partitionConsumer, err := consumer.ConsumePartition(KAFKA_TOPIC, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Failed to create partition consumer for partition %d: %s", partition, err)
				return
			}
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					log.Printf("Error closing partition consumer for partition %d: %s", partition, err)
				}
			}()

			// Process messages
			for message := range partitionConsumer.Messages() {
				log.Printf("Partition %d | Offset %d | Key: %s | Value: %s", message.Partition, message.Offset, string(message.Key), string(message.Value))
			}
		}(partition)
	}

	// Wait for the consumer to finish
	wg.Wait()
}
