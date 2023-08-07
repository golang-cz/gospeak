all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	go generate ./...

run-server:
	@echo "Running server. You can test it with:"
	@echo "$$ go test --serverUrl=http://localhost:8080 ./..."
	@echo
	go run ./cmd/server

test:
	go test ./...
