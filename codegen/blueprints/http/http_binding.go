package http

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type RequestT struct{}
type ResponseT struct{}

func makeHTTPBinding(ctx context.Context, e endpoint.Endpoint, before []httptransport.RequestFunc, after []httptransport.ResponseFunc) http.Handler {
	decode := func(r *http.Request) (interface{}, error) {
		defer r.Body.Close()
		var request RequestT
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, err
		}
		return request, nil
	}
	encode := func(w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
	return httptransport.Server{
		Context:            ctx,
		Endpoint:           e,
		DecodeRequestFunc:  decode,
		EncodeResponseFunc: encode,
		Before:             before,
		After:              append([]httptransport.ResponseFunc{httptransport.SetContentType("application/json; charset=utf-8")}, after...),
	}
}
