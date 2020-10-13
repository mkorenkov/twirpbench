# twirpbench

Twirp is [protobuf](https://developers.google.com/protocol-buffers/docs/proto3)-based service-to-service communication framework, similar to [gRPC](http://www.grpc.io/).
Check out [twirp repo](https://github.com/twitchtv/twirp) to learn more.

Unfrotunately, Twirp has some underlying implementation details that are extremely expensive and there's nothing you can do about it.
Or is it?

```
protoc --proto_path=$GOPATH/src:. --twirp_out=internal --go_out=internal ./internal/rpc/twirpdefault/bloat.proto
protoc --proto_path=$GOPATH/src:. --twirp_out=internal/twirpoptimized --gofast_out=internal/twirpoptimized ./internal/rpc/bloat.proto
```
