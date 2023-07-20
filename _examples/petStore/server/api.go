package server

import (
	"sync"

	"github.com/golang-cz/gospeak/_examples/petStore/proto"
)

type API struct {
	mu       sync.Mutex
	PetStore map[int64]*proto.Pet
	SeqID    int64
}
