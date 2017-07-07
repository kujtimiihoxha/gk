struct FooReply{
}
struct FooRequest{
}

service TestThriftService {
 FooReply Foo (1: FooRequest req)
}

