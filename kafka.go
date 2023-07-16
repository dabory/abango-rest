package abango

import (
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	e "github.com/dabory/abango-rest/etc"
)

var (
	KAFKA_CONN      string
	COMSUMER_TOPICS []string
	KAFKA_TIMEOUT   string
)

// sdfjasldfja
func KafkaInit() {
	KAFKA_CONN = XConfig["KafkaConnString"]
	COMSUMER_TOPICS = strings.Split(strings.Replace(XConfig["ConsumerTopics"], " ", "", -1), ",")
	KAFKA_TIMEOUT = XConfig["KafkaTimeout"]
	e.OkLog("== KAFKA_CONN is : " + KAFKA_CONN + " ==")
	e.OkLog("== COMSUMER_TOPICS is : " + XConfig["ConsumerTopics"] + " ==")
	e.OkLog("== KAFKA_TIMEOUT is : " + KAFKA_TIMEOUT + " ==")
}

func KafkaProducer(key string, headers []*sarama.RecordHeader, message []byte, conCurr string, topic string) (int32, int64, error) {

	kfcf := sarama.NewConfig()
	kfcf.Producer.Retry.Max = 5
	kfcf.Producer.RequiredAcks = sarama.WaitForAll
	kfcf.Producer.Return.Successes = true
	conHeaders := e.ConvertKafkaHeaders(headers)
	// fmt.Println("KAFKA_CONN:", KAFKA_CONN)
	if conCurr == "async" {
		if prd, err := sarama.NewAsyncProducer([]string{KAFKA_CONN}, kfcf); err == nil {
			prd.Input() <- &sarama.ProducerMessage{
				Topic:   topic,
				Key:     sarama.StringEncoder(key),
				Headers: conHeaders,
				Value:   sarama.ByteEncoder(message),
			}
			return 0, 0, nil
		} else {
			return 0, 0, e.MyErr("QEJHDRTTRRW-Kafka-NewAsyncProducer-End", err, true)
		}
	} else if conCurr == "sync" {
		if prd, err := sarama.NewSyncProducer([]string{KAFKA_CONN}, kfcf); err == nil {
			msg := &sarama.ProducerMessage{
				Topic:   topic,
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

	// using for loop
	for i := 0; i < len(COMSUMER_TOPICS); i++ {

		partitions, err := consumer.Partitions(COMSUMER_TOPICS[i])
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
				partitionConsumer, err := consumer.ConsumePartition(COMSUMER_TOPICS[i], partition, sarama.OffsetNewest)
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
}
