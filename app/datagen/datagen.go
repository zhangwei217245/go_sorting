package datagen

import (
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"

	"example.com/kconnect"
)

var (
	idCounter = uint32(1)
)

// RandomAddress returns a random string of length 15-20, with a mixture of numbers, characters, and space, length ranging from 15-20
func RandomAddress() string {
	min := 15
	max := 20
	var length = rand.Intn(max-min+1) + min
	var letters = []rune(" abcdefghijklmnopqrstuvwxyz0123456789")

	return RandomString(length, letters)
}

// RandomName returns a random string of length 10-15, with English characters only.
func RandomName() string {
	min := 10
	max := 15
	var length = rand.Intn(max-min+1) + min
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	return RandomString(length, letters)
}

// RandomString returns a random string of specified length, with characters from specified charset.
func RandomString(length int, charset []rune) string {
	s := make([]rune, length)
	for i := range s {
		s[i] = charset[rand.Intn(len(charset))]
	}
	return string(s)
}

// NewIDString returns a new uint32 ID in its string form.
func NewIDString(channel int) string {
	number := NewID(uint32(channel))
	return strconv.FormatUint(uint64(number), 10)
}

// RandomContinent returns a continent name.
func RandomContinent() string {
	list := []string{"North America", "Asia", "South America", "Europe", "Africa", "Australia"}
	return list[rand.Intn(len(list))]
}

// NewID returns a new ID of uint32. It allows maximum 2^5=32 channels to create ID simultaneously.
func NewID(channel uint32) uint32 {
	i := atomic.AddUint32(&idCounter, 1)
	// make higher 5 bits all zero.
	low27 := i & ((^uint32(0)) >> 5)
	// apply channel id to higher-order 5 bits.
	high5 := channel << 27
	return low27 | high5
}

// NewRow returns a row in the csv data, with all required columns
func NewRow(channel int) string {
	return strings.Join([]string{NewIDString(channel), RandomName(), RandomAddress(), RandomContinent()}, ",")
}

// Generate generates specified number (mcount) of rows using tnum threads and send them to specified kafka topic
func Generate(mcount int64, topic string, kaddr string) {

	producer, err := kconnect.InitProducer(kaddr)
	if err != nil {
		panic(err)
	}
	// logging.DebugLogger.Println("initialized")
	var i int64
	for i = 0; i < mcount; i++ {
		msg := NewRow(0)
		kconnect.Publish(topic, strconv.FormatInt(int64(i), 10), msg, producer)
	}

	if err := producer.Close(); err != nil {
		panic(err)
	}

}
