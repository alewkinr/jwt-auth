package crypto

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
)

type Source struct{}

// New создаем новый источник для псевдо-рандомизации
func New() Source {
	return Source{}
}

// Реалзиуем функции интерфейса rand.Source

// Seed ...
func (s Source) Seed(seed int64) {}

// Int63 ...
func (s Source) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

// Uint64 ...
func (s Source) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
