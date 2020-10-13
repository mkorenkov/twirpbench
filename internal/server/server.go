package server

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const blobSource = "/dev/urandom"

type Request interface {
	GetCompressed() bool
	GetByteSize() int64
	GetKey() string
}

type Response struct {
	Key        string
	Compressed bool
	Value      []byte
}

func GetRandomBytes(req Request) (Response, error) {
	f, err := os.OpenFile(blobSource, os.O_RDONLY, 0755)
	if err != nil {
		return Response{}, errors.Wrapf(err, "Error opening %s", blobSource)
	}
	defer func() {
		if dErr := f.Close(); dErr != nil {
			panic(dErr)
		}
	}()

	randReader := io.LimitReader(f, req.GetByteSize())

	var r io.Reader
	r = randReader

	if req.GetCompressed() {
		pr, pw := io.Pipe()
		r = pr
		go func() {
			gz := gzip.NewWriter(pw)
			defer func() {
				if dErr := gz.Close(); dErr != nil {
					panic(dErr)
				}
				if dErr := pw.Close(); dErr != nil {
					panic(dErr)
				}
			}()
			_, err := io.Copy(pw, randReader)
			if err != nil {
				panic(err)
			}
		}()
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Response{}, errors.Wrapf(err, "Error reading bytes from %s", blobSource)
	}
	return Response{
		Key:        req.GetKey(),
		Compressed: req.GetCompressed(),
		Value:      data,
	}, nil
}
