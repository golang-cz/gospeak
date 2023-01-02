all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	gospeak ./schema/api.go \
	  golang@v0.10.0 -server -pkg=server -out=./server/server.gen.go \
	  golang@v0.10.0 -client -pkg=client -out=./client/client.gen.go \
	  typescript@v0.10.0 -client -out=./client.gen.ts \
	  openapi@v0.10.0 -out=./openapi.gen.yaml

run-server:
	@echo "Running server. You can test it with:"
	@echo "$$ go test --serverUrl=http://localhost:8080 ./..."
	@echo
	go run ./cmd/server

test:
	go test ./...
