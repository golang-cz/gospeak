package main

import (
	"net/http"

	"github.com/golang-cz/go2webrpc/_examples/userStore/server"
)

func main() {
	rpc := &server.RPC{
		UserStore: map[int64]*server.User{},
	}

	apiServer := server.NewUserStoreServer(rpc)
	http.ListenAndServe(":8080", apiServer)
}
