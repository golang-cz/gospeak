build:
	mkdir -p bin
	go build -o bin/gospeak ./cmd/gospeak

generate-petstore: build
	cd _examples/petStore/proto && ../../../bin/gospeak .

clean:
	rm -rf bin

install:
	go install ./cmd/gospeak
