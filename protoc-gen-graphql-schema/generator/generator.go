package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/spec"
	"google.golang.org/protobuf/proto"
)

type Template struct {
	RootPackage *spec.Package

	Packages   []*spec.Package
	Types      []*spec.Message
	Interfaces []*spec.Message
	Enums      []*spec.Enum
	Inputs     []*spec.Message
	Services   []*spec.Service
}

// Generator is struct for analyzing protobuf definition
// and factory graphql definition in protobuf to generate.
type Generator struct {
	files    []*spec.File
	args     *spec.Params
	messages map[string]*spec.Message
	enums    map[string]*spec.Enum
	logger   *Logger
}

func New(files []*spec.File, args *spec.Params) *Generator {

	messages := make(map[string]*spec.Message)
	enums := make(map[string]*spec.Enum)

	for _, f := range files {
		for _, m := range f.Messages() {
			messages[m.FullPath()] = m
		}
		for _, e := range f.Enums() {
			enums[e.FullPath()] = e
		}
	}

	w := io.Discard
	if args.Verbose {
		w = os.Stderr
	}

	return &Generator{
		files:    files,
		args:     args,
		messages: messages,
		enums:    enums,
		logger:   NewLogger(w),
	}
}

func (g *Generator) Generate(tmpl string, fs []string) ([]*plugin.CodeGeneratorResponse_File, error) {

	services, err := g.analyzeServices()
	if err != nil {
		return nil, err
	}

	var outFiles []*plugin.CodeGeneratorResponse_File
	for _, f := range g.files {
		for _, v := range fs {
			if f.Filename() != v {
				continue
			}

			s, ok := services[f.Package()]
			if !ok {
				continue
			}

			// mark as same package definition in file
			g.analyzeEnum(f)
			if err := g.analyzeMessage(f); err != nil {
				return nil, err
			}

			file, err := g.generateFile(f, tmpl, s)
			if err != nil {
				return nil, err
			}
			outFiles = append(outFiles, file)
		}
	}
	return outFiles, nil
}

// nolint: gocognit, funlen, gocyclo
func (g *Generator) generateFile(file *spec.File, tmpl string, services []*spec.Service) (
	*plugin.CodeGeneratorResponse_File,
	error,
) {

	var types, inputs, interfaces []*spec.Message
	var enums []*spec.Enum
	var packages []*spec.Package

	for _, m := range g.messages {
		// skip empty field message, otherwise graphql-go raise error
		if len(m.Fields()) == 0 {
			continue
		}
		if m.IsDepended(spec.DependTypeMessage, file.Package()) {
			switch {
			case file.Package() == m.Package():
				types = append(types, m)
			case spec.IsGooglePackage(m):
				packages = append(packages, spec.NewGooglePackage(m))
			default:
				packages = append(packages, spec.NewPackage(m))
			}
		}
		if m.IsDepended(spec.DependTypeInput, file.Package()) {
			if !spec.IsGooglePackage(m) {
				inputs = append(inputs, m)
			}
		}
		if m.IsDepended(spec.DependTypeInterface, file.Package()) {
			interfaces = append(interfaces, m)
		}
	}

	for _, s := range services {
		for _, q := range s.Queries {
			input, output := q.Input, q.Output
			if input.Package() != file.Package() {
				if spec.IsGooglePackage(input) {
					packages = append(packages, spec.NewGooglePackage(input))
				} else {
					packages = append(packages, spec.NewPackage(input))
				}
			}
			if output.Package() != file.Package() {
				if spec.IsGooglePackage(output) {
					packages = append(packages, spec.NewGooglePackage(output))
				} else {
					packages = append(packages, spec.NewPackage(output))
				}
			}
		}
		for _, m := range s.Mutations {
			input, output := m.Input, m.Output
			if input.Package() != file.Package() {
				if spec.IsGooglePackage(input) {
					packages = append(packages, spec.NewGooglePackage(input))
				} else {
					packages = append(packages, spec.NewPackage(input))
				}
			}
			if output.Package() != file.Package() {
				if spec.IsGooglePackage(output) {
					packages = append(packages, spec.NewGooglePackage(output))
				} else {
					packages = append(packages, spec.NewPackage(output))
				}
			}
		}
	}

	for _, e := range g.enums {
		// skip empty values enum, otherwise graphql-go raise error
		if len(e.Values()) == 0 {
			continue
		}
		if e.IsDepended(spec.DependTypeEnum, file.Package()) {
			if file.Package() == e.Package() || spec.IsGooglePackage(e) {
				enums = append(enums, e)
			} else {
				packages = append(packages, spec.NewPackage(e))
			}
		}
	}

	// drop duplicate packages
	uniquePackages := make([]*spec.Package, 0)
	stack := make(map[string]struct{})
	for _, p := range packages {
		if _, ok := stack[p.Path]; ok {
			continue
		}
		uniquePackages = append(uniquePackages, p)
		stack[p.Path] = struct{}{}
	}

	// Sort by name to avoid to appear some diff on each generation
	sort.Slice(uniquePackages, func(i, j int) bool {
		return uniquePackages[i].Name > uniquePackages[j].Name
	})
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() > types[j].Name()
	})
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name() > enums[j].Name()
	})
	sort.Slice(inputs, func(i, j int) bool {
		return inputs[i].Name() > inputs[j].Name()
	})
	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].Name() > interfaces[j].Name()
	})
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name() > services[j].Name()
	})

	root := spec.NewPackage(file)
	t := &Template{
		RootPackage: root,
		Packages:    uniquePackages,
		Types:       types,
		Enums:       enums,
		Inputs:      inputs,
		Interfaces:  interfaces,
		Services:    services,
	}

	buf := new(bytes.Buffer)
	if tmpl, err := template.New("go").Parse(tmpl); err != nil {
		return nil, err
	} else if err := tmpl.Execute(buf, t); err != nil {
		return nil, err
	}

	out := buf.Bytes()

	// If paths=source_relative option is provided, put generated file relatively
	if g.args.IsSourceRelative() {
		return &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.graphql", root.GeneratedFilenamePrefix)),
			Content: proto.String(string(out)),
		}, nil
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fmt.Sprintf("%s/%s.graphql", root.Path, root.FileName)),
		Content: proto.String(string(out)),
	}, nil
}
