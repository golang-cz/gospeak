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
