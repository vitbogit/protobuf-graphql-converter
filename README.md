# protobuf-graphql-converter

`protobuf-graphql-converter` is a protoc plugin that generates GraphQL schema from Protocol Buffers schema.

## 💪 Current capability

Implemented:

- simple "protobuf message" -> "graphql type" conversion 

- simple "protobuf enum" -> "graphql enum" conversion

## ⬇️ Installation

### TODO: Simple installation using "go get"

### Full installation 

1) Clone repository:
```git clone https://github.com/vitbogit/protobuf-graphql-converter.git```

2) (not recommended, risky) If you want to rebuild whole plugin\`s core, then run:
 ```make plugin```
 This command will rebuild plugin\`s core, defined in include/graphql/graphql.proto file and further stored in graphql/graphql.pb.go file.

 TODO: Windows support. Make commands doesn`t work... 

3) To obtain executable of plugin, run:
 ```make distribute```

Now you should have your plugin\`s executable in dist folder!

 TODO: Windows support. Make commands doesn`t work... 

4) (recommended) Move plugin\`s executable to your $(GOBIN) folder.

TODO: Windows support. Now you need to manually add ".exe" to the file name!

5) (recommended) Move include/graphql/graphql.proto file to your project`s protobuf files folder

## 🚀 Usage

### How to use installed plugin

To use plugin, protoc needs:
- path to plugin\`s executable (if you\`ve done recommended installation steps, then plugin is now in $(GOBIN) so it\`s reachable with PATH, otherwise path to plugin\`s executable should be provided with --plugin=..., when running protoc command)
- path to graphql.proto file (which is at first located in my project at ./include/graphql/graphql.proto, if you\`ve done recommended installation steps it\`s now in one folder with your protobuf files)

If you did steps 4 and 5 of installation guide, you can now simply run:

```
protoc -I. --graphql-schema_out=YOUR_OUT_FOLDER YOUR_PROTO_FILE.proto
```

That will result to generating GraphQL schema out of YOUR_PROTO_FILE.proto

### How to run examples

If everything is installed, fast way to run example is to run `make example` in project\`s root directory or to run `make generate` in examples/greeter directory.