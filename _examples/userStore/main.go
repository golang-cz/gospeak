package main

import (
	"net/http"

	"github.com/golang-cz/go2webrpc/_examples/userStore/rpc"
)

func main() {
	server := &rpc.RPC{}

	http.ListenAndServe(":8080", rpc.NewExampleAPIServer(server))
}
