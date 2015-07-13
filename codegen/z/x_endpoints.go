// Do not edit! Generated by gokit-generate

package z

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)


type XEndpoints struct {
	x X
}

func MakeXEndpoints (x X) XEndpoints {
	return XEndpoints{x}
}


type XClient struct {
	e endpoint.Endpoint
}

func MakeXClient (e endpoint.Endpoint) XClient {
	return XClient{e}
}

func (x XEndpoints) Y (ctx context.Context, request interface{}) (interface{}, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, endpoint.ErrContextCanceled
	}
	req, ok := request.(XYRequest)
	if !ok {
		return nil, endpoint.ErrBadCast
	}
	var err error
	var resp XYResponse
	resp.Int64 = x.x.Y(ctx, req.P, req.Int, req.Int1, req.Int64)
	return resp, err
}

func (x XEndpoints) Z (ctx context.Context, request interface{}) (interface{}, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, endpoint.ErrContextCanceled
	}
	req, ok := request.(XZRequest)
	if !ok {
		return nil, endpoint.ErrBadCast
	}
	var err error
	var resp XZResponse
	resp.R, err = x.x.Z(ctx, req.A, req.B)
	return resp, err
}

// TODO: implement X methods on XClient.
