package ws

import (
	"io"

	"github.com/go-kit/kit/endpoint"

	"context"
)

// Server-Side Codec

// EndpointCodec defines a server Endpoint and its associated codecs
type EndpointCodec struct {
	Endpoint endpoint.Endpoint
	Decode   DecodeRequestFunc
	Encode   EncodeResponseFunc
}

// EndpointCodecMap maps the Request.Method to the proper EndpointCodec
type EndpointCodecMap map[string]EndpointCodec

// DecodeRequestFunc extracts a user-domain io.Reader from the incoming WebSocket Subprotocol
// It's designed to be used in WebSocket servers, for server-side endpoints.
// One straightforward DecodeRequestFunc could be something that unmarshals
// JSON from the request reader to the concrete request type.
type DecodeRequestFunc func(context.Context, io.Reader) (request interface{}, err error)

// EncodeResponseFunc encodes the passed response object to a WebSocket Subprotocol result.
// It's designed to be used in WebSocket servers, for server-side endpoints.
// One straightforward EncodeResponseFunc could be something that JSON encodes
// the object directly.
type EncodeResponseFunc func(context.Context, io.Writer, interface{}) (err error)

// Client-Side Codec

// EncodeRequestFunc encodes the given request object to raw JSON.
// It's designed to be used in JSON RPC clients, for client-side
// endpoints. One straightforward EncodeResponseFunc could be something that
// JSON encodes the object directly.
// type EncodeRequestFunc func(context.Context, interface{}) (response interface{}, err error)

// DecodeResponseFunc extracts a user-domain response object from an JSON RPC
// response object. It's designed to be used in JSON RPC clients, for
// client-side endpoints. It is the responsibility of this function to decide
// whether any error present in the JSON RPC response should be surfaced to the
// client endpoint.
// type DecodeResponseFunc func(context.Context, interface{}) (response interface{}, err error)
