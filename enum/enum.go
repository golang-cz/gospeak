package enum

// type Enum[T EnumType] T
//
// type EnumType interface {
//	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
// }
//
// NOTE: The generic Enum type didn't work, as it failed with:
// "cannot use a type parameter as RHS in type declaration"
// https://github.com/golang/go/issues/45639

type Int int
type Uint uint

type Int8 int8
type Uint8 uint8

type Int16 int16
type Uint16 uint16

type Int32 int32
type Uint32 uint32

type Int64 int64
type Uint64 uint64

// webrpc TODO: string ENUM
// https://github.com/webrpc/webrpc/issues/203
