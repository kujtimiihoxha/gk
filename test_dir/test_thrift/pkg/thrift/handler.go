package thrift

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/kujtimiihoxha/gk/test_dir/test_thrift/pkg/endpoints"
	thriftTestThrift "github.com/kujtimiihoxha/gk/test_dir/test_thrift/pkg/thrift/gen-go/test_thrift"
)

type thriftServer struct {
	ctx context.Context
	foo endpoint.Endpoint
}

// MakeThriftHandler makes a set of endpoints available as a thrift server.
func MakeThriftHandler(ctx context.Context, endpoints endpoints.Endpoints) (req thriftTestThrift.TestThriftService) {
	req = &thriftServer{
		ctx: ctx,
		foo: endpoints.FooEndpoint,
	}
	return req
}

// DecodeThriftFooRequest is a func that converts a
// thrift request to a user-domain request. Primarily useful in a server.
// TODO: Do not forget to implement the decoder.
func DecodeThriftFooRequest(r *thriftTestThrift.FooRequest) (req endpoints.FooRequest, err error) {
	err = errors.New("'Foo' Decoder is not impelement")
	return req, err
}

// EncodeThriftFooResponse is a func that converts a
// user-domain response to a thrift reply. Primarily useful in a server.
// TODO: Do not forget to implement the encoder.
func EncodeThriftFooResponse(reply interface{}) (rep thriftTestThrift.FooReply, err error) {
	err = errors.New("'Foo' Encoder is not impelement")
	return rep, err
}

func (s *thriftServer) Foo(req *thriftTestThrift.FooRequest) (rep *thriftTestThrift.FooReply, err error) {
	request, err := DecodeThriftFooRequest(req)
	if err != nil {
		return nil, err
	}
	response, err := s.foo(s.ctx, request)
	if err != nil {
		return nil, err
	}
	r, err := EncodeThriftFooResponse(response)
	rep = &r
	return rep, err
}
