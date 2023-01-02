package server

import (
	"context"
)

func (s *API) GetPet(ctx context.Context, ID int64) (pet *Pet, err error) {
	pet, ok := s.PetStore[ID]
	if !ok {
		return nil, ErrorNotFound("pet(%v) not found", ID)
	}
	return pet, nil
}

func (s *API) ListPets(ctx context.Context) (pets []*Pet, err error) {
	pets = make([]*Pet, 0, len(s.PetStore))
	for _, pet := range s.PetStore {
		pets = append(pets, pet)
	}

	return pets, nil
}

func (s *API) CreatePet(ctx context.Context, pet *Pet) (*Pet, error) {
	if pet == nil {
		return nil, ErrorInvalidArgument("pet", "pet is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	pet.ID = s.SeqID
	s.SeqID++

	s.PetStore[pet.ID] = pet
	return pet, nil
}

func (s *API) UpdatePet(ctx context.Context, ID int64, pet *Pet) (*Pet, error) {
	if pet == nil {
		return nil, ErrorInvalidArgument("pet", "pet is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.PetStore[pet.ID]
	if !ok {
		return nil, ErrorNotFound("pet(%v) not found", ID)
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
