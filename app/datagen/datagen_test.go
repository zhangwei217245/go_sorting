package datagen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	Generate(100000000, "test", "127.0.0.1:9092")
}
