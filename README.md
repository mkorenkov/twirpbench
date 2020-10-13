# twirpbench

benchmark for twirp

```
protoc --proto_path=$GOPATH/src:. --twirp_out=internal --go_out=internal ./internal/rpc/bloat.proto
#protoc --proto_path=$GOPATH/src:. --gofast_out=internal/io ./internal/rpc/bloat.proto
```
