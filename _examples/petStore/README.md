# GoSpeak userStore example

1. Read and update API definition at [proto/api.go](./proto/api.go)
2. Generate client/server code
    ```
    make generate
    ```
3. Run server
    ```
    $ make run-server
    ```
4. Run client tests
    ```
    $ go test ./...
    ```
