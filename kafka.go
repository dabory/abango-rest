package abango

import (
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

func KafkaConsumer(ConsumeHandler func(msg *sarama.ConsumerMessage), topic string) {

	// Create a new configuration for the consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Specify the list of brokers
	brokers := []string{KAFKA_CONN}

	// Create a new consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		e.OkLog("Failed to create consumer of topic : " + topic + " == : " + err.Error())
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			e.OkLog("Error closing consumer: of topic : " + topic + " == : " + err.Error())
		}
	}()

	// Consume messages from each partition asynchronously
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		e.OkLog("Failed to get partitions: " + topic + " == : " + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()

			// Create a new partition consumer
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				e.OkLog("Failed to create partition for consumer: " + topic + " partition: " + e.NumToStr(partition) + " " + err.Error())
				return
			}
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					e.OkLog("Error closing partition for consumer: " + topic + " partition: " + e.NumToStr(partition) + " " + err.Error())
				}
			}()

			// Process messages
			for msg := range partitionConsumer.Messages() {
				ConsumeHandler(msg)
				e.OkLog("Consuming topic: " + topic + " partition: " + e.NumToStr(partition))
				// log.Printf("Partition-kk %d | Offset %d | Key: %s | Value: %s", message.Partition, message.Offset, string(message.Key), string(message.Value))
			}
		}(partition)
	}
	// Wait for the consumer to finish
	wg.Wait()
	// }
}
