package schema

import "context"

//go:generate gospeak -schema=./api.go -interface=UserStore -out=./webrpc.gen.json
//go:generate webrpc-gen -schema=./webrpc.gen.json -target=golang -server -pkg=server -out=../server/server.gen.go
//go:generate webrpc-gen -schema=./webrpc.gen.json -target=golang -client -pkg=client -out=../client/client.gen.go

type UserStore interface {
	UpsertUser(ctx context.Context, new *User) (user *User, err error)
	GetUser(ctx context.Context, ID int64) (user *User, err error)
	ListUsers(ctx context.Context) (users []*User, err error)
	DeleteUser(ctx context.Context, ID int64) error
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
