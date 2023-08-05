package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-cz/gospeak/_examples/petStore/proto"
	"github.com/golang-cz/gospeak/_examples/petStore/server"
)

func main() {
	api := &server.API{
		PetStore: map[int64]*proto.Pet{},
		SeqID:    1,
	}

	r := chi.NewRouter()
	r.Use(DebugPayload)
	r.Handle("/*", proto.NewPetStoreServer(api))

	log.Println("Serving PetStore API at :8080")
	http.ListenAndServe(":8080", r)
}

func DebugPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bufReq, bufResp bytes.Buffer

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &bufReq))
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		rw.Tee(&bufResp)

		next.ServeHTTP(rw, r)

		log.Println("request:", bufReq.String())
		log.Println("response:", bufResp.String())
	})
}
