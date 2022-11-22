package server

type RPC struct {
	UserStore    map[int64]*User
	SequentialID int64
}
