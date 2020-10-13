package bench

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"

	"github.com/mkorenkov/twirpbench/internal/twirpdefault"
	twdefault "github.com/mkorenkov/twirpbench/internal/twirpdefault/rpc/bloat"
	twoptimized "github.com/mkorenkov/twirpbench/internal/twirpoptimized/rpc/bloat"
)

type request struct {
	Key        string
	Compressed bool
	ByteSize   int64
}

type client interface {
	MakeRequest(ctx context.Context, url string, req request) error
}

type twDefaultClient struct{}

func (*twDefaultClient) MakeRequest(ctx context.Context, url string, req request) error {
	client := twdefault.NewBloatProtobufClient(url, http.DefaultClient)
	res, err := client.GetBlob(ctx, &twdefault.BlobRequest{
		Key:        req.Key,
		Compressed: req.Compressed,
		ByteSize:   req.ByteSize,
	})
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("result is nil")
	}
	return nil
}

type twOptimizedClient struct{}

func (*twOptimizedClient) MakeRequest(ctx context.Context, url string, req request) error {
	client := twoptimized.NewBloatProtobufClient(url, http.DefaultClient)
	res, err := client.GetBlob(ctx, &twoptimized.BlobRequest{
		Key:        req.Key,
		Compressed: req.Compressed,
		ByteSize:   req.ByteSize,
	})
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("result is nil")
	}
	return nil
}

func BenchmarkUnscopedSearch(b *testing.B) {
	ctx, cancelFn := context.WithTimeout(context.TODO(), 10*time.Minute)
	defer cancelFn()

	defaultTwirpHandler := twdefault.NewBloatServer(&twirpdefault.Server{})
	defaultTwirpClient := &twDefaultClient{}

	testCases := []struct {
		name         string
		compressed   bool
		size         int64
		twirpHandler http.Handler
		twirpClient  client
	}{
		{"raw-300K", false, 300000, defaultTwirpHandler, defaultTwirpClient},
		{"raw-1M", false, 1000000, defaultTwirpHandler, defaultTwirpClient},
		{"raw-10M", false, 10000000, defaultTwirpHandler, defaultTwirpClient},
	}

	for _, testCase := range testCases {
		b.Run(testCase.name, func(b *testing.B) {
			ts := httptest.NewServer(testCase.twirpHandler)
			defer ts.Close()

			req := request{
				Key:        testCase.name,
				Compressed: testCase.compressed,
				ByteSize:   testCase.size,
			}

			var start, end runtime.MemStats
			runtime.GC()

			runtime.ReadMemStats(&start)
			b.SetBytes(testCase.size)
			b.ResetTimer()

			if err := testCase.twirpClient.MakeRequest(ctx, ts.URL, req); err != nil {
				log.Fatal(err)
			}

			b.ReportAllocs()
			runtime.ReadMemStats(&end)
			b.ReportMetric(float64(end.TotalAlloc-start.TotalAlloc)/1024/1024, "TotalAlloc(MiB)")
		})
		b.StopTimer()
	}
}
