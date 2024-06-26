package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/generator"
	"github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/spec"
	"github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/templates"
	"google.golang.org/protobuf/proto"
)

var version = "dev"
var printVersion = flag.Bool("v", false, "show binary version")

func main() {
	flag.Parse()
	if *printVersion {
		io.WriteString(os.Stdout, version)
		os.Exit(0)
	}

	var genError error

	// To add support of optional fields (protobuf version 3),
	// we shoudld enable plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL.
	// In bit mask it`s just 1.
	// Read more:
	// https://pkg.go.dev/google.golang.org/protobuf/types/pluginpb#CodeGeneratorResponse
	var responseFeaturesMask uint64 = 1

	resp := &plugin.CodeGeneratorResponse{SupportedFeatures: &responseFeaturesMask}
	defer func() {
		// If some error has been occurred in generate process,
		// add error message to plugin response
		if genError != nil {
			message := genError.Error()
			resp.Error = &message
		}
		buf, err := proto.Marshal(resp)
		if err != nil {
			log.Fatalln(err)
		}
		os.Stdout.Write(buf)
	}()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(os.Stdin); err != nil {
		genError = err
		return
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(buf.Bytes(), &req); err != nil {
		genError = err
		return
	}

	var parameter string
	if req.Parameter != nil {
		parameter = req.GetParameter()
	}
	args, err := spec.NewParams(parameter)
	if err != nil {
		genError = err
		return
	}

	// We're dealing with each descriptors to out wrapper struct
	// in order to access easily plugin options, package name, comment, etc...
	var files []*spec.File
	for _, f := range req.GetProtoFile() {
		files = append(files, spec.NewFile(f, req.GetCompilerVersion(), args.FieldCamelCase))
	}

	g := generator.New(files, args)
	var ftg []string
	for _, f := range req.GetFileToGenerate() {
		if !args.IsExclude(f) {
			ftg = append(ftg, f)
		}
	}
	if len(ftg) > 0 {
		genFiles, err := g.Generate(templates.GraphqlTemplate, ftg)
		if err != nil {
			genError = err
			return
		}
		resp.File = append(resp.File, genFiles...)
	}
}
