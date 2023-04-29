package hack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"google.golang.org/protobuf/encoding/protojson"
	kyaml "sigs.k8s.io/yaml"
)

func TestCB(t *testing.T) {
	f, err := os.ReadFile("cloudbuild.yaml")
	if err != nil {
		t.Fatal(err)
	}
	js, err := kyaml.YAMLToJSON(f)
	if err != nil {
		t.Fatal(err)
	}
	cb := &cloudbuildpb.Build{}
	err = protojson.Unmarshal(js, cb)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(cb); err != nil {
		t.Fatal(err)
	}
	// kyaml.JSONToYAML()
	// for _, s := range cb.Steps {
	// 	fmt.Println(protojson.Format(s))
	// 	// sh.Parse() // s.Script
	// }

	// fmt.Println(cb.String())
	fmt.Println("=====================================")

	bla := &cloudbuildpb.Build{
		Name:             "",
		Id:               "",
		ProjectId:        "",
		Status:           0,
		StatusDetail:     "",
		Source:           nil,
		Steps:            nil,
		Results:          nil,
		CreateTime:       nil,
		StartTime:        nil,
		FinishTime:       nil,
		Timeout:          nil,
		Images:           nil,
		QueueTtl:         nil,
		Artifacts:        nil,
		LogsBucket:       "",
		SourceProvenance: nil,
		BuildTriggerId:   "",
		Options:          nil,
		LogUrl:           "",
		Substitutions:    nil,
		Tags:             nil,
		Secrets:          nil,
		Timing:           nil,
		Approval:         nil,
		ServiceAccount:   "",
		AvailableSecrets: nil,
		Warnings:         nil,
		FailureInfo:      nil,
	}

	// desc.LoadMessageDescriptorForMessage(cb.ProtoReflect())
	// p := protoprint.Printer{}
	fmt.Println(cb.ProtoReflect().Descriptor())
	// fmt.Println(p.PrintProtoToString())
	// fmt.Println(
	// 	valast.StringWithOptions(
	// 		cb.ProtoMessage, &valast.Options{
	// 			Unqualify: false,
	// 			// PackagePath:  "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb",
	// 			// PackageName:  "cloudbuildpb",
	// 			ExportedOnly: true,
	// 			// PackagePathToName: "cb",
	// 		},
	// 	),
	// )
	// spew.Dump(cb) // &cloudbuildpb.Build{...}
	// cloudbuildpb.Build
}

func TestExportedFunctions(t *testing.T) {
	// write a function that parsea .go file and list all the exported functions
	// and their comments
	// there are many go packages that can help with this
	// for example https://pkg.go.dev/golang.org/x/tools/go/packages
	// and https://pkg.go.dev/golang.org/x/tools/go/ast/astutil
	//
	// the output should be a map[string]string
	// where the key is the function name and the value is the comment
	// for example:
	// map[string]string{
	// 	"main": "main is the entry point for the application",
	// 	"run": "run is the main logic for the application",
	// }

	// hint: you can use the go/ast package to parse the .go file
	// and the go/doc package to extract the comments
	// https://pkg.go.dev/golang.org/x/tools/go/ast/astutil
	// https://pkg.go.dev/golang.org/x/tools/go/doc

}
