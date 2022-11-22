package server

import "context"

func (s *RPC) UpsertUser(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, ErrorInvalidArgument("user", "user is required")
	}

	if user.ID == 0 {
		s.SequentialID += 1
		user.ID = s.SequentialID
	}

	s.UserStore[user.ID] = user
	return user, nil
}

func (s *RPC) GetUser(ctx context.Context, ID int64) (user *User, err error) {
	user, ok := s.UserStore[ID]
	if !ok {
		return nil, ErrorNotFound("user(%v) not found", ID)
	}
	return user, nil
}

func (s *RPC) ListUsers(ctx context.Context) (users []*User, err error) {
	for _, user := range s.UserStore {
		users = append(users, user)
	}
	return users, nil
}

func (s *RPC) DeleteUser(ctx context.Context, ID int64) error {
	delete(s.UserStore, ID)
	return nil
}
