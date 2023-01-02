package server

import (
	"sync"
)

type API struct {
	mu       sync.Mutex
	PetStore map[int64]*Pet
	SeqID    int64
}
