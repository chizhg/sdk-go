// +build kafka

package kafka_sarama_test

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/binding"
	"github.com/cloudevents/sdk-go/binding/test"
	"github.com/cloudevents/sdk-go/bindings/kafka_sarama"
)

// Avoid DCE
var ProducerMessage *sarama.ProducerMessage

var (
	ctxSkipKey                  context.Context
	ctx                         context.Context
	eventWithoutKey             cloudevents.Event
	eventWithKey                cloudevents.Event
	structuredMessageWithoutKey binding.Message
	structuredMessageWithKey    binding.Message
	binaryMessageWithoutKey     binding.Message
	binaryMessageWithKey        binding.Message
)

func init() {
	ctxSkipKey = kafka_sarama.WithSkipKeyExtension(context.TODO())
	ctx = context.TODO()

	eventWithoutKey = test.FullEvent()
	eventWithKey = test.FullEvent()
	eventWithKey.SetExtension("key", "aaaaaa")

	structuredMessageWithoutKey = test.NewMockStructuredMessage(eventWithoutKey)
	structuredMessageWithKey = test.NewMockStructuredMessage(eventWithKey)
	binaryMessageWithoutKey = test.NewMockBinaryMessage(eventWithoutKey)
	binaryMessageWithKey = test.NewMockBinaryMessage(eventWithKey)
}

func BenchmarkEncodeStructuredMessageSkipKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProducerMessage = &sarama.ProducerMessage{}
		Err = kafka_sarama.EncodeKafkaProducerMessage(ctxSkipKey, structuredMessageWithoutKey, ProducerMessage, binding.TransformerFactories{})
	}
}

func BenchmarkEncodeStructuredMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProducerMessage = &sarama.ProducerMessage{}
		Err = kafka_sarama.EncodeKafkaProducerMessage(ctx, structuredMessageWithKey, ProducerMessage, binding.TransformerFactories{})
	}
}

func BenchmarkEncodeBinaryMessageSkipKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProducerMessage = &sarama.ProducerMessage{}
		Err = kafka_sarama.EncodeKafkaProducerMessage(ctxSkipKey, binaryMessageWithoutKey, ProducerMessage, binding.TransformerFactories{})
	}
}

func BenchmarkEncodeBinaryMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProducerMessage = &sarama.ProducerMessage{}
		Err = kafka_sarama.EncodeKafkaProducerMessage(ctx, binaryMessageWithKey, ProducerMessage, binding.TransformerFactories{})
	}
}
