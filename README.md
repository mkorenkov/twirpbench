# twirpbench

Twirp is [protobuf](https://developers.google.com/protocol-buffers/docs/proto3)-based service-to-service communication framework, similar to [gRPC](http://www.grpc.io/).
Check out [twirp repo](https://github.com/twitchtv/twirp) to learn more.

This repo was created in attempt to improve protobuf twirp client. Replacing `ioutilReadAll` gives about 2x better Total Allocations and 20% in latency by getting rid of unnecessary `ioutil.ReadAll(resp.Body)` allocations. Unfortunately, golang protobuf implementation is not `io.Reader` friendly, so this performance win requires some jumps through the hoop.

Benchmark results
```
BenchmarkTwirp/twirp-raw-300K-16         	1000000000	         0.00163 ns/op	184269977629.62 MB/s	         2.64 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-raw-1M-16           	1000000000	         0.00392 ns/op	255041533513.73 MB/s	         5.98 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-raw-10M-16          	1000000000	         0.0426 ns/op	234744373940.28 MB/s	        83.1 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-raw-100M-16         	1000000000	         0.387 ns/op	258257712225.45 MB/s	       703 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-gz-300K-16          	1000000000	         0.00237 ns/op	126818632024.33 MB/s	         3.44 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-gz-1M-16            	1000000000	         0.00442 ns/op	226198643577.21 MB/s	         6.79 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-gz-10M-16           	1000000000	         0.0397 ns/op	251609362676.25 MB/s	        83.9 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/twirp-gz-100M-16          	1000000000	         0.339 ns/op	294597201047.31 MB/s	       704 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-raw-300K-16      	1000000000	         0.00147 ns/op	203911429031.69 MB/s	         1.34 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-raw-1M-16        	1000000000	         0.00308 ns/op	324346134410.34 MB/s	         3.02 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-raw-10M-16       	1000000000	         0.0326 ns/op	306301247644.62 MB/s	        41.6 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-raw-100M-16      	1000000000	         0.312 ns/op	320623666896.89 MB/s	       351 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-gz-300K-16       	1000000000	         0.00156 ns/op	192274906970.99 MB/s	         2.14 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-gz-1M-16         	1000000000	         0.00378 ns/op	264353048783.71 MB/s	         3.83 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-gz-10M-16        	1000000000	         0.0304 ns/op	328668929617.98 MB/s	        42.4 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkTwirp/maxtwirp-gz-100M-16       	1000000000	         0.267 ns/op	374342562278.78 MB/s	       352 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
```

## How to get these results for your project?

### Prerequisites
```
brew install protobuf
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
# get vanilla twirp protoc plugin
go get -u github.com/twitchtv/twirp/protoc-gen-twirp
# get minimal gogo/protobuf protoc plugin
go get -u github.com/mkorenkov/twirpbench/cmd/protoc-gen-maxgo

# install protoc plugin with ioutil.ReadAll fix
# it's a hacky way due to `internal` module usage of `twitchtv/twirp` protoc plugin
rm -rf "$GOPATH/src/github.com/twitchtv/twirp"
mkdir -p "$GOPATH/src/github.com/twitchtv/twirp"
git clone "https://github.com/mkorenkov/twirp.git" "$GOPATH/src/github.com/twitchtv/twirp"
cd "$GOPATH/src/github.com/twitchtv/twirp"
go build -o "$GOPATH/bin/protoc-gen-maxtwirp" ./protoc-gen-twirp
rm -rf "$GOPATH/src/github.com/twitchtv/twirp"
```

### check default twirp generated code
```
protoc --proto_path=$GOPATH/src:. --twirp_out=internal --go_out=internal ./internal/rpc/twirpdefault/bloat.proto
```

### get ioutil.ReadAll optimized protobuf code
optimized twirp generated code (with `mkorenkov/twirp` protoc plugin):
```
protoc --proto_path=$GOPATH/src:. --maxtwirp_out=internal --maxgo_out=internal ./internal/rpc/twirpoptimized/bloat.proto
```

NOTE: this is equivallent to manually applying `0001-patch.patch` to the generated file.
```
# `github.com/twitchtv/twirp` plugin will be used instead of `github.com/mkorenkov/twirp`
protoc --proto_path=$GOPATH/src:. --twirp_out=internal --maxgo_out=internal ./internal/rpc/twirpoptimized/bloat.proto
patch -p1 < 0001-patch.patch
```
