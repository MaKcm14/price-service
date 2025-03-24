package kafka

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/MaKcm14/price-service/internal/controller/chttp"
	"github.com/MaKcm14/price-service/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// Producer defines the logic of kafka's producing.
type Producer struct {
	producer sarama.SyncProducer
	logger   *slog.Logger
}

func NewProducer(log *slog.Logger, brokers []string) (Producer, error) {
	const op = "kafka.new-producer"

	conf := sarama.NewConfig()

	conf.Producer.Timeout = 2 * time.Second
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, conf)

	if err != nil {
		log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return Producer{}, fmt.Errorf("error of the %s: %s", op, err)
	}

	return Producer{
		producer: producer,
		logger:   log,
	}, nil
}

// SendProductsMessage sends the products response to the client.
func (p Producer) SendProductsMessage(products []entities.ProductSample, request dto.ProductRequest) {
	const op = "kafka.send-products-message"

	response := chttp.NewProductResponse(products)
	buf, _ := json.Marshal(response)

	recordHeaders := make([]sarama.RecordHeader, 0, len(request.Headers))

	for key, val := range request.Headers {
		recordHeaders = append(recordHeaders, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(val),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   productsTopicName,
		Value:   sarama.ByteEncoder(buf),
		Headers: recordHeaders,
	}

	count := 0
	for _, _, err := p.producer.SendMessage(msg); err != nil && count != 5; count++ {
		p.logger.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		time.Sleep(time.Millisecond * 50)
	}
}

// Close shuts down the producer and releases another resources.
func (p Producer) Close() {
	p.producer.Close()
}
