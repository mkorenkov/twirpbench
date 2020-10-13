package twirpoptimized

import (
	"context"

	"github.com/mkorenkov/twirpbench/internal/server"
	"github.com/mkorenkov/twirpbench/internal/twirpoptimized/rpc/bloat"
)

type Server struct{}

func (*Server) GetBlob(ctx context.Context, req *bloat.BlobRequest) (*bloat.Blob, error) {
	res, err := server.GetRandomBytes(req)
	if err != nil {
		return nil, err
	}
	return &bloat.Blob{
		Key:        res.Key,
		Compressed: res.Compressed,
		Value:      res.Value,
	}, nil
}
