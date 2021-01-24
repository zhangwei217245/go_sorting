package kconnect

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"example.com/logging"
	"github.com/Shopify/sarama"
)

// InitProducer returns a initialized SyncProducer instance and an error instance
func InitProducer(kafkaConn string) (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = logging.DebugLogger // log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

// Publish sends message with key through the SyncProducer instance to a specified kafka topic.
func Publish(topic string, key string, message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		logging.ErrorLogger.Println("Error publish: ", err.Error())
	}

	logging.InfoLogger.Printf("Sending to topic %-12v Partition: %-3v Offset: %-10v key: %-10v msg: %s\n", topic, p, o, key, message)
}

// InitConsumer returns a initialized Consumer instance and an error instance
func InitConsumer(kafkaConn string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = "testClient"

	// Specify brokers address. This is default one
	brokers := []string{kafkaConn}

	// Create new consumer
	return sarama.NewConsumer(brokers, config)
}

// Consume reads the specified number of messages from specified topic, and process the message with procMsg function closure.
func Consume(topic string, procMsg func(int32, int64, string, string), consumer sarama.Consumer) {

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		panic(err)
	}

	var (
		messages = make(chan *sarama.ConsumerMessage, 256)
		closing  = make(chan struct{})
		wg       sync.WaitGroup
	)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-signals
		logging.InfoLogger.Println("Initiating shutdown of consumer...")
		close(closing)
	}()

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}

		go func(pc sarama.PartitionConsumer) {
			<-closing
			pc.AsyncClose()
		}(pc)

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				messages <- msg
			}
		}(pc)
	}

	go func() {
		for msg := range messages {
			procMsg(msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		}
	}()

	wg.Wait()
	logging.InfoLogger.Println("Done consuming topic", topic)
	close(messages)

}
