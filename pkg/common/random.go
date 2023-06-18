package common

import (
	"bytes"
	"math/rand"
)

func RandomToken(rng *rand.Rand, size int) string {
	hex := []byte{
		'0', '1', '2', '3',
		'4', '5', '6', '7',
		'8', '9', 'a', 'b',
		'c', 'd', 'e', 'f',
	}

	var buffer bytes.Buffer

	for i := 0; i < size; i++ {
		index := rng.Intn(len(hex))
		buffer.WriteByte(hex[index])
	}

	return buffer.String()
}
