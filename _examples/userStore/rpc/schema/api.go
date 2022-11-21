package schema

import "context"

//go:generate go2webrpc -schema=./api.go -interface=ExampleAPI -out=./webrpc.json
//go:generate webrpc-gen -schema=./webrpc.json -target=golang -server -pkg=rpc -out=../server.gen.go
//go:generate webrpc-gen -schema=./webrpc.json -target=golang -client -pkg=users -out=../../pkg/users/client.gen.go

type ExampleAPI interface {
	GetSession(ctx context.Context) (user *User, err error)

	Get(ctx context.Context, ID int64) (user *User, err error)
	ListUsers(ctx context.Context) (users []*User, err error)
}

type User struct {
	ID   int64
	UID  string
	Name string
}

// PostA *Post `db:"-" json:"postA"`
// PostX XXX   `db:"-" json:"postX"`
// PostY YYY   `db:"-" json:"postY"`
// PostZ ZZZ   `db:"-" json:"postZ"`

// type XXX []Post
// type YYY *Post

// type F struct{}
// type ZZZ []F
