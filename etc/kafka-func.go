// Author : Eric Kim
// Build Date : 6 Jul 2023  Last Update 02 Aug 2018
// All rights are reserved.

package etc

import (
	"bytes"

	"github.com/Shopify/sarama"
)

func KafkaHeaderValue(headers []*sarama.RecordHeader, key []byte) string {
	for _, header := range headers {
		if bytes.Equal(header.Key, key) {
			return string(header.Value)
		}
	}
	return ""
}

func ConvertKafkaHeaders(headers []*sarama.RecordHeader) []sarama.RecordHeader {
	convertedHeaders := make([]sarama.RecordHeader, len(headers))
	for i, header := range headers {
		convertedHeaders[i] = *header
	}
	return convertedHeaders
}
