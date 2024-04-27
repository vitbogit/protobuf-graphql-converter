package generator

import "github.com/vitbogit/protobuf-graphql-converter/protoc-gen-graphql-schema/spec"

func (g *Generator) getMessage(name string) *spec.Message {
	if v, ok := g.messages[name]; ok {
		return v
	} else if v, ok := g.messages["."+name]; ok {
		return v
	}
	return nil
}

func (g *Generator) getEnum(name string) *spec.Enum {
	if v, ok := g.enums[name]; ok {
		return v
	} else if v, ok := g.enums["."+name]; ok {
		return v
	}
	return nil
}
