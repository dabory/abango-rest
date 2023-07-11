package abango

import (
	"log"
	"sync"

	"github.com/Shopify/sarama"
	e "github.com/dabory/abango-rest/etc"
)

var (
	KAFKA_CONN    string
	KAFKA_TOPIC   string
	KAFKA_TIMEOUT string
)

// sdfjasldfja
func KafkaInit() {
	KAFKA_CONN = XConfig["KafkaConnString"]
	KAFKA_TOPIC = XConfig["KafkaTopic"]
	KAFKA_TIMEOUT = XConfig["KafkaTimeout"]
	e.OkLog("== KAFKA_CONN is : " + KAFKA_CONN + " ==")
	e.OkLog("== KAFKA_TOPIC is : " + KAFKA_TOPIC + " ==")
	e.OkLog("== KAFKA_TIMEOUT is : " + KAFKA_TIMEOUT + " ==")
}

func KafkaProducer(key string, headers []*sarama.RecordHeader, message []byte, conCurr string) (int32, int64, error) {

	kfcf := sarama.NewConfig()
	kfcf.Producer.Retry.Max = 5
	kfcf.Producer.RequiredAcks = sarama.WaitForAll
	kfcf.Producer.Return.Successes = true
	conHeaders := e.ConvertKafkaHeaders(headers)

	if conCurr == "async" {
		if prd, err := sarama.NewAsyncProducer([]string{KAFKA_CONN}, kfcf); err == nil {
			prd.Input() <- &sarama.ProducerMessage{
				Topic:   KAFKA_TOPIC,
				Key:     sarama.StringEncoder(key),
				Headers: conHeaders,
				Value:   sarama.ByteEncoder(message),
			}
			return 0, 0, nil
		} else {
			return 0, 0, e.MyErr("QEJHDRTTRRW-Kafka-NewSyncProducer-End", err, true)
		}
	} else if conCurr == "sync" {
		if prd, err := sarama.NewSyncProducer([]string{KAFKA_CONN}, kfcf); err == nil {
			msg := &sarama.ProducerMessage{
				Topic:   KAFKA_TOPIC,
				Key:     sarama.StringEncoder(key),
				Headers: conHeaders,
				Value:   sarama.ByteEncoder(message),
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

func KafkaConsumer(ConsumeHandler func(msg *sarama.ConsumerMessage)) {

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
			for msg := range partitionConsumer.Messages() {
				ConsumeHandler(msg)
				// log.Printf("Partition-kk %d | Offset %d | Key: %s | Value: %s", message.Partition, message.Offset, string(message.Key), string(message.Value))
			}
		}(partition)
	}

	// Wait for the consumer to finish
	wg.Wait()
}
