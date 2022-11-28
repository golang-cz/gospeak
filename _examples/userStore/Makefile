all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	gospeak ./schema/api.go \
	  golang -server -pkg=server -out=./server/server.gen.go \
	  golang -client -pkg=client -out=./client/client.gen.go \
	  typescript -client -out=./client.ts.go \
	  openapi -out=./openapi.yaml

run-server:
	go run ./cmd/server

run-client-tests:
	go test ./...
