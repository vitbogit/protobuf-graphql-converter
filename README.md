# protobuf-graphql-converter

`protobuf-graphql-converter` is a protoc plugin that generates GraphQL schema from Protocol Buffers schema.

## Current capability

Implemented:

- simple "protobuf message" -> "graphql type" conversion 

- simple "protobuf enum" -> "graphql enum" conversion

## Installation

Plugin depends on new changes made to protobuf and graphql, therefore
there are several installations methods:

### TODO: Simple installation using "go get"

### Full installation 

1) Clone repository:
```git clone https://github.com/vitbogit/protobuf-graphql-converter.git```

2) (not recomended, risky) If you want to rebuild whole plugin`s core, then run:
 ```make plugin```
 This command will rebuild plugin`s core, defined in ./include/graphql/graphql.proto file.

3) To obtain executable of plugin, run:
 ```make distribute```

Now you should have your plugin`s executable in ./dist folder!

4) (recomended) Move plugin`s executable to your $(GOBIN) folder.
TODO: Windows support. Now you need to manually add ".exe" to the file name!

5) (recomended) Move ./include/graphql/graphql.proto to your protobuf files folder

## Usage

To use plugin, protoc needs "to know":
- path to plugin`s executable (plugin might just be in go/bin)
- path to graphql.proto file (which is at first located in my project at ./include/graphql/graphql.proto)

If you did steps 4 and 5 from installation guide, you can now simply run:

```
protoc -I. --graphql-schema_out=YOUR_OUT_FOLDER YOUR_PROTO_FILE.proto
```

## Example

If everything is installed, fast way to run example is to run `make example` in project`s root directory or to run `make generate` in examples/greeter directory.