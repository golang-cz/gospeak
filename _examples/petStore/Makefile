all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	gospeak ./schema/api.go \
	  golang@v0.8.0 -server -pkg=server -out=./server/server.gen.go \
	  golang@v0.8.0 -client -pkg=client -out=./client/client.gen.go \
	  typescript@v0.7.0 -client -out=./client.ts.go \
	  openapi@v0.7.0 -out=./openapi.yaml

run-server:
	go run ./cmd/server

run-client-tests:
	go test ./...
