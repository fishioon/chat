package chat

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

/*
func newSessionID() string {
	var sid [16]byte
	for i := 0; i < len(sid)/4; i++ {
		r := rand.Uint32()
		for j := 0; j < 4 && i*4+j < len(sid); j++ {
			sid[i*4+j] = byte(r)
			r >>= 8
		}
	}
	return string(sid[:])
}
*/
