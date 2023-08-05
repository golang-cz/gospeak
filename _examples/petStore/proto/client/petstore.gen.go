// PetStore vTODO 204f6b26587305ef3a4c043b8636035ada3889ef
// --
// Code generated by webrpc-gen@v0.13.0-dev with github.com/webrpc/gen-golang@tags/opts_types generator. DO NOT EDIT.
//
// gospeak .
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
	"time"

	"github.com/google/uuid"
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
	return "204f6b26587305ef3a4c043b8636035ada3889ef"
}

//
// Types
//

type Status int

const (
	Status_approved Status = 0
	Status_pending Status = 1
	Status_closed Status = 2
	Status_new Status = 3
)

var Status_name = map[int]string{
	0: "approved",
	1: "pending",
	2: "closed",
	3: "new",
}

var Status_value = map[string]int{
	"approved": 0,
	"pending": 1,
	"closed": 2,
	"new": 3,
}

func (x Status) String() string {
	return Status_name[int(x)]
}

func (x Status) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(Status_name[int(x)])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

func (x *Status) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*x = Status(Status_value[j])
	return nil
}

type Tag struct {
	ID int64 `json:"ID"`
	Name string `json:"Name"`
}

type Pet struct {
	ID int64 `json:"id,string"`
	UUID uuid.UUID `json:"uuid,string"`
	Name string `json:"name"`
	Available bool `json:"available"`
	PhotoURLs []string `json:"photoUrls"`
	Tags []Tag `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Tag Tag `json:"Tag"`
	TagPtr *Tag `json:"TagPtr"`
	TagsPtr []Tag `json:"TagsPtr"`
	Status Status `json:"status"`
}

type PetStore interface {
	CreatePet(ctx context.Context, new *Pet) (*Pet, error)
	DeletePet(ctx context.Context, ID int64) error
	GetPet(ctx context.Context, ID int64) (*Pet, error)
	ListPets(ctx context.Context) ([]*Pet, error)
	UpdatePet(ctx context.Context, ID int64, update *Pet) (*Pet, error)
}

var WebRPCServices = map[string][]string{
	"PetStore": {
		"CreatePet",
		"DeletePet",
		"GetPet",
		"ListPets",
		"UpdatePet",
	},
}

//
// Client
//

const PetStorePathPrefix = "/rpc/PetStore/"

type petStoreClient struct {
	client HTTPClient
	urls	 [5]string
}

func NewPetStoreClient(addr string, client HTTPClient) PetStore {
	prefix := urlBase(addr) + PetStorePathPrefix
	urls := [5]string{
		prefix + "CreatePet",
		prefix + "DeletePet",
		prefix + "GetPet",
		prefix + "ListPets",
		prefix + "UpdatePet",
	}
	return &petStoreClient{
		client: client,
		urls:	 urls,
	}
}

func (c *petStoreClient) CreatePet(ctx context.Context, new *Pet) (*Pet, error) {
	in := struct {
		Arg0 *Pet `json:"new"`
	}{new}
	out := struct {
		Ret0 *Pet `json:"pet"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[0], in, &out)
	return out.Ret0, err
}

func (c *petStoreClient) DeletePet(ctx context.Context, ID int64) error {
	in := struct {
		Arg0 int64 `json:"ID"`
	}{ID}

	err := doJSONRequest(ctx, c.client, c.urls[1], in, nil)
	return err
}

func (c *petStoreClient) GetPet(ctx context.Context, ID int64) (*Pet, error) {
	in := struct {
		Arg0 int64 `json:"ID"`
	}{ID}
	out := struct {
		Ret0 *Pet `json:"pet"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[2], in, &out)
	return out.Ret0, err
}

func (c *petStoreClient) ListPets(ctx context.Context) ([]*Pet, error) {
	out := struct {
		Ret0 []*Pet `json:"pets"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[3], nil, &out)
	return out.Ret0, err
}

func (c *petStoreClient) UpdatePet(ctx context.Context, ID int64, update *Pet) (*Pet, error) {
	in := struct {
		Arg0 int64 `json:"ID"`
		Arg1 *Pet `json:"update"`
	}{ID, update}
	out := struct {
		Ret0 *Pet `json:"pet"`
	}{}

	err := doJSONRequest(ctx, c.client, c.urls[4], in, &out)
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
		return ErrorWithCause(ErrWebrpcRequestFailed, fmt.Errorf("failed to marshal JSON body: %w", err))
	}
	if err = ctx.Err(); err != nil {
		return ErrorWithCause(ErrWebrpcRequestFailed, fmt.Errorf("aborted because context was done: %w", err))
	}

	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), "application/json")
	if err != nil {
		return ErrorWithCause(ErrWebrpcRequestFailed, fmt.Errorf("could not build request: %w", err))
	}
	resp, err := client.Do(req)
	if err != nil {
		return ErrorWithCause(ErrWebrpcRequestFailed, err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = ErrorWithCause(ErrWebrpcRequestFailed, fmt.Errorf("failed to close response body: %w", cerr))
		}
	}()

	if err = ctx.Err(); err != nil {
		return ErrorWithCause(ErrWebrpcRequestFailed, fmt.Errorf("aborted because context was done: %w", err))
	}

	if resp.StatusCode != 200 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to read server error response body: %w", err))
		}

		var rpcErr WebRPCError
		if err := json.Unmarshal(respBody, &rpcErr); err != nil {
			return ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to unmarshal server error: %w", err))
		}
		if rpcErr.Cause != "" {
			rpcErr.cause = errors.New(rpcErr.Cause)
		}
		return rpcErr
	}

	if out != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to read response body: %w", err))
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return ErrorWithCause(ErrWebrpcBadResponse, fmt.Errorf("failed to unmarshal JSON response body: %w", err))
		}
	}

	return nil
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

//
// Errors
//

type WebRPCError struct {
	Name       string `json:"error"`
	Code       int    `json:"code"`
	Message    string `json:"msg"`
	Cause      string `json:"cause,omitempty"`
	HTTPStatus int    `json:"status"`
	cause      error
}

var _ error = WebRPCError{}

func (e WebRPCError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s %d: %s: %v", e.Name, e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s %d: %s", e.Name, e.Code, e.Message)
}

func (e WebRPCError) Is(target error) bool {
	if rpcErr, ok := target.(WebRPCError); ok {
		return rpcErr.Code == e.Code
	}
	return errors.Is(e.cause, target)
}

func (e WebRPCError) Unwrap() error {
	return e.cause
}

func ErrorWithCause(rpcErr WebRPCError, cause error) WebRPCError {
	err := rpcErr
	err.cause = cause
	err.Cause = cause.Error()
	return err
}

// Webrpc errors
var (
	ErrWebrpcEndpoint = WebRPCError{Code: 0, Name: "WebrpcEndpoint", Message: "endpoint error", HTTPStatus: 400}
	ErrWebrpcRequestFailed = WebRPCError{Code: -1, Name: "WebrpcRequestFailed", Message: "request failed", HTTPStatus: 400}
	ErrWebrpcBadRoute = WebRPCError{Code: -2, Name: "WebrpcBadRoute", Message: "bad route", HTTPStatus: 404}
	ErrWebrpcBadMethod = WebRPCError{Code: -3, Name: "WebrpcBadMethod", Message: "bad method", HTTPStatus: 405}
	ErrWebrpcBadRequest = WebRPCError{Code: -4, Name: "WebrpcBadRequest", Message: "bad request", HTTPStatus: 400}
	ErrWebrpcBadResponse = WebRPCError{Code: -5, Name: "WebrpcBadResponse", Message: "bad response", HTTPStatus: 500}
	ErrWebrpcServerPanic = WebRPCError{Code: -6, Name: "WebrpcServerPanic", Message: "server panic", HTTPStatus: 500}
	ErrWebrpcInternalError = WebRPCError{Code: -7, Name: "WebrpcInternalError", Message: "internal error", HTTPStatus: 500}
)
