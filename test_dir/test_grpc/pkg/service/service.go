package service

import "context"

// Implement yor service methods methods.
// e.x: Foo(ctx context.Context,s string)(rs string, err error)
type TestGrpcService interface {
	Foo(ctx context.Context, s string) (rs string, err error)
}

type stubTestGrpcService struct{}

// Get a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New() (s *stubTestGrpcService) {
	s = &stubTestGrpcService{}
	return s
}

// Implement the business logic of Foo
func (te *stubTestGrpcService) Foo(ctx context.Context, s string) (rs string, err error) {
	return rs, err
}
