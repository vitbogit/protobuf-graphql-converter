# protobuf-graphql-converter

`protobuf-graphql-converter` is a protoc plugin that generates graphql schema from Protocol Buffers schema

## Current capability

Implemented:

- simple "protobuf message" -> "graphql type" conversion 

- simple "protobuf enum" -> "graphql enum" conversion

## Installation

Plugin depends on new changes made to protobuf and graphql, therefore
there are several installations methods:

### TODO: Simple installation using "go get"

### Another simple installation (not rebuilding plugin`s core)

1) Clone repository:
```git clone git clone https://github.com/vitbogit/protobuf-graphql-converter.git```

2) Run with bash:
 ```make distribute```

Now you should have your plugin`s executable in ./dist folder!

3) TODO: install executable to go/bin folder

### Full installation 

1) Clone repository:
```git clone git clone https://github.com/vitbogit/protobuf-graphql-converter.git```

2) Run with bash:
 ```make plugin```
 That will rebuild plugin`s core from ./include/graphql/graphql.proto file.

3) If you are lucky (that means, if the plugin`s logic is compatible with current protobuf
and graphql standards), then this command will work:
 ```make distribute```

Now you should have your plugin`s executable in ./dist folder!

4) TODO: install executable to go/bin folder

## Usage

To use plugin, protoc needs to know:
- path to plugin`s executable (plugin might just be in go/bin)
- path to graphql.proto file (which is at first located in my project at ./include/graphql/graphql.proto)

For simplicity, you can just copy graphql.proto file to your protobufs schemas.
In case you have moved graphql.proto file to your protobuf and not installed plugin to go/bin, you can copy plugin`s executable to same folder as protobuf files and then just run:

```
protoc -I. --plugin=./protoc-gen-graphql-schema --graphql-schema_out=./ YOUR_PROTO_NAME.proto
```

## Example

If everything is installed, just run `make example` in project`s root directory