package rpc

import "context"

func (s *RPC) GetSession(ctx context.Context) (user *User, err error) {
	return s.userStore[0], nil
}
