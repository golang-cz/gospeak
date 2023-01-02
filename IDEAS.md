## Enums in Go
```
import github.com/golang-cz/gospeak

type Status = gospeak.Enum[int64, string]{
  0: "unknown",
  1: "draft",
  2: "scheduled",
  3: "published",
  4: "deleted",
}
```

## Test recursive types

// PostA *Post `db:"-" json:"postA"`
// PostX XXX   `db:"-" json:"postX"`
// PostY YYY   `db:"-" json:"postY"`
// PostZ ZZZ   `db:"-" json:"postZ"`

// type XXX []Post
// type YYY *Post

// type F struct{}
// type ZZZ []F


## YAML configuration file?

```yaml
userStoreApi:
  schema: api.yml
  interfaces: []
  gen:
    - golang@v0.8.0 -server -pkg=server -out=./server/server.gen.go
    - golang@v0.8.0 -client -pkg=client -out=./client/client.gen.go
    - typescript@v0.7.0 -client -out=./client.ts.go
    - openapi@v0.7.0 -out=./openapi.yaml
```