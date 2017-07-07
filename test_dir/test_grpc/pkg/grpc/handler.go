package grpc

import (
	"context"
	"errors"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/kujtimiihoxha/gk/test_dir/test_grpc/pkg/endpoints"
	"github.com/kujtimiihoxha/gk/test_dir/test_grpc/pkg/grpc/pb"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	foo grpctransport.Handler
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints endpoints.Endpoints) (req pb.TestGrpcServer) {
	req = &grpcServer{
		foo: grpctransport.NewServer(
			endpoints.FooEndpoint,
			DecodeGRPCFooRequest,
			EncodeGRPCFooResponse,
		),
	}
	return req
}

// DecodeGRPCFooRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
// TODO: Do not forget to implement the decoder, you can find an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/transport_grpc.go#L62-L65
func DecodeGRPCFooRequest(_ context.Context, grpcReq interface{}) (req interface{}, err error) {
	err = errors.New("'Foo' Decoder is not impelement")
	return req, err
}

// EncodeGRPCFooResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
// TODO: Do not forget to implement the encoder, you can find an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/transport_grpc.go#L62-L65
func EncodeGRPCFooResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	err = errors.New("'Foo' Encoder is not impelement")
	return res, err
}

func (s *grpcServer) Foo(ctx oldcontext.Context, req *pb.FooRequest) (rep *pb.FooReply, err error) {
	_, rp, err := s.foo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	rep = rp.(*pb.FooReply)
	return rep, err
}
