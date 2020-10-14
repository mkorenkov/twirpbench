// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gofast/main.go

package main

import (
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()

	vanity.ForEachFile(files, vanity.TurnOffGogoImport)

	//vanity.ForEachFile(files, vanity.TurnOnMarshalerAll)
	vanity.ForEachFile(files, vanity.TurnOnSizerAll)
	//vanity.ForEachFile(files, vanity.TurnOnUnmarshalerAll)

	resp := command.Generate(req)
	command.Write(resp)
}
