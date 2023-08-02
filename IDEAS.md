## Test recursive types

// PostA *Post `db:"-" json:"postA"`
// PostX XXX   `db:"-" json:"postX"`
// PostY YYY   `db:"-" json:"postY"`
// PostZ ZZZ   `db:"-" json:"postZ"`

// type XXX []Post
// type YYY *Post

// type F struct{}
// type ZZZ []F
