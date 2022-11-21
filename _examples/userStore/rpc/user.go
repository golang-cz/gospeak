package rpc

import "context"

func (s *RPC) Get(ctx context.Context, ID int64) (user *User, err error) {
	user, ok := s.userStore[ID]
	if !ok {
		return nil, ErrorNotFound("user(%v) not found", ID)
	}
	return user, nil
}

func (s *RPC) ListUsers(ctx context.Context) (users []*User, err error) {
	for _, user := range s.userStore {
		users = append(users, user)
	}
	return users, nil
}
