package main

import (
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()

	vanity.ForEachFile(files, vanity.TurnOffGogoImport)
	// no special features needed:
	// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gofast/main.go for comparison

	resp := command.Generate(req)
	command.Write(resp)
}
