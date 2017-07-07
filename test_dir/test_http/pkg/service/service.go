package service

import "context"

// Implement yor service methods methods.
// e.x: Foo(ctx context.Context,s string)(rs string, err error)
type TestHttpService interface {
	Foo(ctx context.Context, s string) (rs string, err error)
}

type stubTestHttpService struct{}

// Get a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New() (s *stubTestHttpService) {
	s = &stubTestHttpService{}
	return s
}

// Implement the business logic of Foo
func (te *stubTestHttpService) Foo(ctx context.Context, s string) (rs string, err error) {
	return rs, err
}
