package generator

import (
	"errors"

	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/vitbogit/protobuf-graphql-converter/graphql"
	"github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/spec"
)

// nolint: interfacer
func (g *Generator) analyzeMessage(file *spec.File) error {
	for _, m := range g.messages {
		if m.Package() != file.Package() {
			continue
		}
		m.Depend(spec.DependTypeMessage, file.Package())
		m.Depend(spec.DependTypeInput, file.Package())
		if err := g.analyzeFields(file.Package(), m, m.Fields(), false, false); err != nil {
			return err
		}
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeEnum(file *spec.File) {
	for _, e := range g.enums {
		if e.Package() != file.Package() {
			continue
		}
		e.Depend(spec.DependTypeEnum, file.Package())
	}
}

func (g *Generator) analyzeServices() (map[string][]*spec.Service, error) {
	services := make(map[string][]*spec.Service)

	for _, f := range g.files {
		services[f.Package()] = []*spec.Service{}

		for _, s := range f.Services() {
			if err := g.analyzeService(f, s); err != nil {
				return nil, err
			}
			if len(s.Queries) > 0 || len(s.Mutations) > 0 {
				services[f.Package()] = append(services[f.Package()], s)
			}
		}
	}
	return services, nil
}

func (g *Generator) analyzeService(f *spec.File, s *spec.Service) error {
	for _, m := range s.Methods() {
		if m.Schema == nil {
			continue
		}
		var input, output *spec.Message

		if input = g.getMessage(m.Input()); input == nil {
			return errors.New("failed to resolve input message: " + m.Input())
		}
		if output = g.getMessage(m.Output()); output == nil {
			return errors.New("failed to resolve output message: " + m.Output())
		}

		switch m.Schema.GetType() {
		case graphql.GraphqlType_QUERY, graphql.GraphqlType_RESOLVER:
			q := spec.NewQuery(m, input, output, g.args.FieldCamelCase)
			if err := g.analyzeQuery(f, q); err != nil {
				return err
			}
			s.Queries = append(s.Queries, q)
		case graphql.GraphqlType_MUTATION:
			mu := spec.NewMutation(m, input, output, g.args.FieldCamelCase)
			if err := g.analyzeMutation(f, mu); err != nil {
				return err
			}
			s.Mutations = append(s.Mutations, mu)
		}
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeQuery(f *spec.File, q *spec.Query) error {
	g.logger.Write("package %s depends on query request %s", f.Package(), q.Input.FullPath())
	q.Input.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Input, q.PluckRequest(), false, false); err != nil {
		return err
	}

	q.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), q.Output, q.PluckResponse(), false, false); err != nil {
		return err
	}
	return nil
}

// nolint: interfacer
func (g *Generator) analyzeMutation(f *spec.File, m *spec.Mutation) error {
	g.logger.Write("package %s depends on mutation request %s", f.Package(), m.Input.FullPath())
	m.Input.Depend(spec.DependTypeInput, f.Package())
	if err := g.analyzeFields(f.Package(), m.Input, m.PluckRequest(), true, false); err != nil {
		return err
	}
	m.Output.Depend(spec.DependTypeMessage, f.Package())
	if err := g.analyzeFields(f.Package(), m.Output, m.PluckResponse(), false, false); err != nil {
		return err
	}
	return nil
}

func (g *Generator) analyzeFields(
	rootPkg string,
	orig *spec.Message,
	fields []*spec.Field,
	asInput,
	recursive bool,
) error {

	for _, f := range fields {
		switch f.Type() {
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			m := g.getMessage(f.TypeName())
			if m == nil {
				return errors.New("failed to resolve field message type: " + f.TypeName())
			}
			f.DependType = m
			if asInput {
				g.logger.Write("package %s depends on input %s", rootPkg, m.FullPath())
				m.Depend(spec.DependTypeInput, rootPkg)
			} else {
				g.logger.Write("package %s depends on message %s", rootPkg, m.FullPath())
				switch {
				case m == orig:
					g.logger.Write("%s has cyclic dependencies of field %s\n", m.Name(), f.Name())
					f.IsCyclic = true
					m.Depend(spec.DependTypeInterface, rootPkg)
				case !recursive:
					m.Depend(spec.DependTypeMessage, rootPkg)
				default:
					return nil
				}
			}

			// Guard from recursive with infinite loop
			if m != orig {
				if err := g.analyzeFields(rootPkg, m, m.Fields(), asInput, true); err != nil {
					return err
				}
			}
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			e := g.getEnum(f.TypeName())
			if e == nil {
				return errors.New("failed to resolve field enum name: " + f.TypeName())
			}
			f.DependType = e
			g.logger.Write("package %s depends on enum %s", rootPkg, e.FullPath())
			e.Depend(spec.DependTypeEnum, rootPkg)
		}
	}
	return nil
}
