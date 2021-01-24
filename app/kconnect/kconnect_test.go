package kconnect

import (
	"fmt"
	"testing"
)

func TestPublish(t *testing.T) {
	// for n := 0; n < 4; n++ {
	// go func() {
	// 	producer, err := InitProducer("127.0.0.1:9092")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for i := 0; i < 1000; i++ {
	// 		msg := strconv.FormatInt(int64(i), 10)
	// 		Publish("test", msg, msg, producer)
	// 	}
	// }()
	// }

}

func procMsg(partition int32, cursor int64, key string, msg string) {
	fmt.Println("key:", key, ",msg:", msg, ", partition:", partition, ", offset:", cursor)
}
func TestConsume(t *testing.T) {
	consumer, err := InitConsumer("127.0.0.1:9092")
	if err != nil {
		panic(err)
	}
	Consume("source", procMsg, consumer)

	if err := consumer.Close(); err != nil {
		panic(nil)
	}
}
