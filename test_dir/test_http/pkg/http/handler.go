package http

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/kujtimiihoxha/gk/test_dir/test_http/pkg/endpoints"
)

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/foo", httptransport.NewServer(
		endpoints.FooEndpoint,
		DecodeFooRequest,
		EncodeFooResponse,
	))
	return m
}

// DecodeFooRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded sum request from the HTTP request body. Primarily useful in a server.
func DecodeFooRequest(_ context.Context, r *http.Request) (req interface{}, err error) {
	req = endpoints.FooRequest{}
	err = json.NewDecoder(r.Body).Decode(&r)
	return req, err
}

// EncodeFooResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeFooResponse(_ context.Context, w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	json.NewEncoder(w).Encode(response)
	return err
}
