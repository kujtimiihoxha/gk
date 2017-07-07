package service

import "context"

// Implement yor service methods methods.
// e.x: Foo(ctx context.Context,s string)(rs string, err error)
type TestThriftService interface {
	Foo(ctx context.Context, s string) (rs string, err error)
}

type stubTestThriftService struct{}

// Get a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New() (s *stubTestThriftService) {
	s = &stubTestThriftService{}
	return s
}

// Implement the business logic of Foo
func (te *stubTestThriftService) Foo(ctx context.Context, s string) (rs string, err error) {
	return rs, err
}
