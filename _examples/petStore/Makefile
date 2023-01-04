all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	gospeak ./schema/api.go \
	  golang@v0.10.0 -server -pkg=server -out=./server/server.gen.go \
	  golang@v0.10.0 -client -pkg=client -out=./client/petstore.gen.go \
	  typescript@v0.10.0 -client -out=./petstore.gen.ts \
	  openapi@v0.10.0 -out=./petstore.gen.yaml

run-server:
	@echo "Running server. You can test it with:"
	@echo "$$ go test --serverUrl=http://localhost:8080 ./..."
	@echo
	go run ./cmd/server

test:
	go test ./...
