package server

import (
	"context"

	"github.com/golang-cz/gospeak/_examples/petStore/proto"
)

func (s *API) GetPet(ctx context.Context, ID int64) (pet *proto.Pet, err error) {
	pet, ok := s.PetStore[ID]
	if !ok {
		return nil, proto.ErrPetNotFound.WithCausef("pet id(%v) not found", ID)
	}
	return pet, nil
}

func (s *API) ListPets(ctx context.Context) (pets []*proto.Pet, err error) {
	pets = make([]*proto.Pet, 0, len(s.PetStore))
	for _, pet := range s.PetStore {
		pets = append(pets, pet)
	}

	return pets, nil
}

func (s *API) CreatePet(ctx context.Context, pet *proto.Pet) (*proto.Pet, error) {
	if pet == nil {
		return nil, proto.ErrInvalidRequest.WithCausef("pet is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	pet.ID = s.SeqID
	s.SeqID++

	s.PetStore[pet.ID] = pet
	return pet, nil
}

func (s *API) UpdatePet(ctx context.Context, ID int64, pet *proto.Pet) (*proto.Pet, error) {
	if pet == nil {
		return nil, proto.ErrInvalidRequest.WithCausef("pet is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.PetStore[pet.ID]
	if !ok {
		return nil, proto.ErrPetNotFound.WithCausef("pet id(%v) not found", ID)
	}

	s.PetStore[pet.ID] = pet
	return pet, nil
}

func (s *API) DeletePet(ctx context.Context, ID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.PetStore, ID)
	return nil
}
