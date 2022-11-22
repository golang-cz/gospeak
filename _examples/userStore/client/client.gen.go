// UserStore vTODO 4b69dce4a22b51dd15700c6efabf661ecedcddc9
// --
// Code generated by webrpc-gen@0.7.6 with golang generator. DO NOT EDIT.
//
// webrpc-gen -schema=./webrpc.gen.json -target=golang -client -pkg=client -out=../client/client.gen.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "vTODO"
}

// Schema hash generated from your RIDL schema
func WebRPCSchemaHash() string {
	return "4b69dce4a22b51dd15700c6efabf661ecedcddc9"
}

//
// Types
//

type User struct {
	ID   int64  `json:"ID"`
	UID  string `json:"UID"`
	Name string `json:"Name"`
}

type UserStore interface {
	DeleteUser(ctx context.Context, ID int64) error
	GetUser(ctx context.Context, ID int64) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	UpsertUser(ctx context.Context, new *User) (*User, error)
}

var WebRPCServices = map[string][]string{
	"UserStore": {
		"DeleteUser",
		"GetUser",
		"ListUsers",
		"UpsertUser",
	},
}

//
// Client
//

const UserStorePathPrefix = "/rpc/UserStore/"

type userStoreClient struct {
	client HTTPClient
	urls   [4]string
}

func NewUserStoreClient(addr string, client HTTPClient) UserStore {
	prefix := urlBase(addr) + UserStorePathPrefix
	urls := [4]string{
		prefix + "DeleteUser",
		prefix + "GetUser",
		prefix + "ListUsers",
		prefix + "UpsertUser",
	}
	return &userStoreClient{
		client: client,
		urls:   urls,
	}
}

func (c *userStoreClient) DeleteUser(ctx context.Context, ID int64) error {
	in := struct {
		Arg0 int64 `json:"iD"`
	}{ID}

	err := doJSONRequest(ctx, c.client, c.urls[0], in, nil)
	return err
}

func (c *userStoreClient) GetUser(ctx context.Context, ID int64) (*User, error) {
	in := struct {
		Arg0 int64 `json:"iD"`
	}{ID}
	out := struct {
		Ret0 *User `json:"user"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[1], in, &out)
	return out.Ret0, err
}

func (c *userStoreClient) ListUsers(ctx context.Context) ([]*User, error) {
	out := struct {
		Ret0 []*User `json:"users"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[2], nil, &out)
	return out.Ret0, err
}

func (c *userStoreClient) UpsertUser(ctx context.Context, new *User) (*User, error) {
	in := struct {
		Arg0 *User `json:"new"`
	}{new}
	out := struct {
		Ret0 *User `json:"user"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[3], in, &out)
	return out.Ret0, err
}

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	if headers, ok := HTTPRequestHeaders(ctx); ok {
		for k := range headers {
			for _, v := range headers[k] {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}

// doJSONRequest is common code to make a request to the remote service.
func doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) error {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return clientError("failed to marshal json request", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("request failed", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	if out != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return clientError("failed to read response body", err)
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return clientError("failed to unmarshal json response body", err)
		}
		if err = ctx.Err(); err != nil {
			return clientError("aborted because context was done", err)
		}
	}

	return nil
}

// errorFromResponse builds a webrpc Error from a non-200 HTTP response.
func errorFromResponse(resp *http.Response) Error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read server error response body", err)
	}

	var respErr ErrorPayload
	if err := json.Unmarshal(respBody, &respErr); err != nil {
		return clientError("failed unmarshal error response", err)
	}

	errCode := ErrorCode(respErr.Code)

	if HTTPStatusFromErrorCode(errCode) == 0 {
		return ErrorInternal("invalid code returned from server error response: %s", respErr.Code)
	}

	return &rpcErr{
		code:  errCode,
		msg:   respErr.Msg,
		cause: errors.New(respErr.Cause),
	}
}

func clientError(desc string, err error) Error {
	return WrapError(ErrInternal, err, desc)
}

func WithHTTPRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {
	if _, ok := h["Accept"]; ok {
		return nil, errors.New("provided header cannot set Accept")
	}
	if _, ok := h["Content-Type"]; ok {
		return nil, errors.New("provided header cannot set Content-Type")
	}

	copied := make(http.Header, len(h))
	for k, vv := range h {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}

	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil
}

func HTTPRequestHeaders(ctx context.Context) (http.Header, bool) {
	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)
	return h, ok
}

//
// Helpers
//

type ErrorPayload struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Cause  string `json:"cause,omitempty"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}

type Error interface {
	// Code is of the valid error codes
	Code() ErrorCode

	// Msg returns a human-readable, unstructured messages describing the error
	Msg() string

	// Cause is reason for the error
	Cause() error

	// Error returns a string of the form "webrpc error <Code>: <Msg>"
	Error() string

	// Error response payload
	Payload() ErrorPayload
}

func Errorf(code ErrorCode, msgf string, args ...interface{}) Error {
	msg := fmt.Sprintf(msgf, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code)}
}

func WrapError(code ErrorCode, cause error, format string, args ...interface{}) Error {
	msg := fmt.Sprintf(format, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg, cause: cause}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code), cause: cause}
}

func Failf(format string, args ...interface{}) Error {
	return Errorf(ErrFail, format, args...)
}

func WrapFailf(cause error, format string, args ...interface{}) Error {
	return WrapError(ErrFail, cause, format, args...)
}

func ErrorNotFound(format string, args ...interface{}) Error {
	return Errorf(ErrNotFound, format, args...)
}

func ErrorInvalidArgument(argument string, validationMsg string) Error {
	return Errorf(ErrInvalidArgument, argument+" "+validationMsg)
}

func ErrorRequiredArgument(argument string) Error {
	return ErrorInvalidArgument(argument, "is required")
}

func ErrorInternal(format string, args ...interface{}) Error {
	return Errorf(ErrInternal, format, args...)
}

type ErrorCode string

const (
	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	ErrUnknown ErrorCode = "unknown"

	// Fail error. General failure error type.
	ErrFail ErrorCode = "fail"

	// Canceled indicates the operation was cancelled (typically by the caller).
	ErrCanceled ErrorCode = "canceled"

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	ErrInvalidArgument ErrorCode = "invalid argument"

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	ErrDeadlineExceeded ErrorCode = "deadline exceeded"

	// NotFound means some requested entity was not found.
	ErrNotFound ErrorCode = "not found"

	// BadRoute means that the requested URL path wasn't routable to a webrpc
	// service and method. This is returned by the generated server, and usually
	// shouldn't be returned by applications. Instead, applications should use
	// NotFound or Unimplemented.
	ErrBadRoute ErrorCode = "bad route"

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	ErrAlreadyExists ErrorCode = "already exists"

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	ErrPermissionDenied ErrorCode = "permission denied"

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	ErrUnauthenticated ErrorCode = "unauthenticated"

	// ResourceExhausted indicates some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	ErrResourceExhausted ErrorCode = "resource exhausted"

	// FailedPrecondition indicates operation was rejected because the system is
	// not in a state required for the operation's execution. For example, doing
	// an rmdir operation on a directory that is non-empty, or on a non-directory
	// object, or when having conflicting read-modify-write on the same resource.
	ErrFailedPrecondition ErrorCode = "failed precondition"

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	ErrAborted ErrorCode = "aborted"

	// OutOfRange means operation was attempted past the valid range. For example,
	// seeking or reading past end of a paginated collection.
	//
	// Unlike InvalidArgument, this error indicates a problem that may be fixed if
	// the system state changes (i.e. adding more items to the collection).
	//
	// There is a fair bit of overlap between FailedPrecondition and OutOfRange.
	// We recommend using OutOfRange (the more specific error) when it applies so
	// that callers who are iterating through a space can easily look for an
	// OutOfRange error to detect when they are done.
	ErrOutOfRange ErrorCode = "out of range"

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	ErrUnimplemented ErrorCode = "unimplemented"

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	ErrInternal ErrorCode = "internal"

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	ErrUnavailable ErrorCode = "unavailable"

	// DataLoss indicates unrecoverable data loss or corruption.
	ErrDataLoss ErrorCode = "data loss"

	// ErrNone is the zero-value, is considered an empty error and should not be
	// used.
	ErrNone ErrorCode = ""
)

func HTTPStatusFromErrorCode(code ErrorCode) int {
	switch code {
	case ErrCanceled:
		return 408 // RequestTimeout
	case ErrUnknown:
		return 400 // Bad Request
	case ErrFail:
		return 422 // Unprocessable Entity
	case ErrInvalidArgument:
		return 400 // BadRequest
	case ErrDeadlineExceeded:
		return 408 // RequestTimeout
	case ErrNotFound:
		return 404 // Not Found
	case ErrBadRoute:
		return 404 // Not Found
	case ErrAlreadyExists:
		return 409 // Conflict
	case ErrPermissionDenied:
		return 403 // Forbidden
	case ErrUnauthenticated:
		return 401 // Unauthorized
	case ErrResourceExhausted:
		return 403 // Forbidden
	case ErrFailedPrecondition:
		return 412 // Precondition Failed
	case ErrAborted:
		return 409 // Conflict
	case ErrOutOfRange:
		return 400 // Bad Request
	case ErrUnimplemented:
		return 501 // Not Implemented
	case ErrInternal:
		return 500 // Internal Server Error
	case ErrUnavailable:
		return 503 // Service Unavailable
	case ErrDataLoss:
		return 500 // Internal Server Error
	case ErrNone:
		return 200 // OK
	default:
		return 0 // Invalid!
	}
}

func IsErrorCode(err error, code ErrorCode) bool {
	if rpcErr, ok := err.(Error); ok {
		if rpcErr.Code() == code {
			return true
		}
	}
	return false
}

func IsValidErrorCode(code ErrorCode) bool {
	return HTTPStatusFromErrorCode(code) != 0
}

type rpcErr struct {
	code  ErrorCode
	msg   string
	cause error
}

func (e *rpcErr) Code() ErrorCode {
	return e.code
}

func (e *rpcErr) Msg() string {
	return e.msg
}

func (e *rpcErr) Cause() error {
	return e.cause
}

func (e *rpcErr) Error() string {
	if e.cause != nil && e.cause.Error() != "" {
		if e.msg != "" {
			return fmt.Sprintf("webrpc %s error: %s -- %s", e.code, e.cause.Error(), e.msg)
		} else {
			return fmt.Sprintf("webrpc %s error: %s", e.code, e.cause.Error())
		}
	} else {
		return fmt.Sprintf("webrpc %s error: %s", e.code, e.msg)
	}
}

func (e *rpcErr) Payload() ErrorPayload {
	statusCode := HTTPStatusFromErrorCode(e.Code())
	errPayload := ErrorPayload{
		Status: statusCode,
		Code:   string(e.Code()),
		Msg:    e.Msg(),
		Error:  e.Error(),
	}
	if e.Cause() != nil {
		errPayload.Cause = e.Cause().Error()
	}
	return errPayload
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	// For Client
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}

	// For Server
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"}

	HTTPRequestCtxKey = &contextKey{"HTTPRequest"}

	ServiceNameCtxKey = &contextKey{"ServiceName"}

	MethodNameCtxKey = &contextKey{"MethodName"}
)