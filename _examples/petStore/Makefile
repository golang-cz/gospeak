all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

proto: generate

generate:
	go generate ./...

run:
	@echo "Running server. You can test it with:"
	@echo "$$ go test --serverUrl=http://localhost:8080 ./..."
	@echo
	go run ./cmd/server

test:
	go test ./...

swagger:
	echo "http://localhost:8088"
	docker run -p 8088:8080 -v $$(pwd)/proto/docs:/docs -e SWAGGER_JSON=/docs/petstore.gen.yaml swaggerapi/swagger-ui
