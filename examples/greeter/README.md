# Greeter example

Simple example of generating GraphQL schema out of protobuf files.

You can run example directly with protoc or using make.

Run example with protoc:
```
protoc -I. -I../../include/graphql --plugin=../../dist/protoc-gen-graphql-schema --graphql-schema_out=./ greeter.proto
```

Run example with make:
```
make generate
```