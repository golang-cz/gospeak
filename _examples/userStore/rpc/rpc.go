package rpc

type RPC struct {
	userStore map[int64]*User
}

var Server = &RPC{
	userStore: map[int64]*User{
		1: &User{ID: 1, UID: "golang-cz", Name: "Golang.cz"},
	},
}
