# protobuf-graphql-converter

`protobuf-graphql-converter` is a protoc plugin that generates GraphQL schema from Protocol Buffers schema.

❗This project is a fork, the architecture of the original generator has not changed much. Also the project had insufficient git commits.

❗This project is no longer maintained.

## 💪 Current capability

Implemented:

- simple "protobuf message" -> "graphql type" conversion, including:

    - lower_snake_case field names to camelCase field names conversion;
    - optional fields handling (generator\`s rule is: resulting field in GraphQL\`s type is not marked non-nullable only if it\`s marked as optional in initial proto file)

- simple "protobuf enum" -> "graphql enum" conversion

Planning:

- conversion improvements 
- better english in README

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

That will result to generating GraphQL schema out of YOUR_PROTO_FILE.proto. Also make sure that YOUR_OUT_FOLDER exists (or simply provide "./").

### How to run examples

If everything is installed, fast way to run example is to run `make example` in project\`s root directory OR to run `make generate` in examples/greeter directory. That will result to generating examples/greeter/greeter.graphql from examples/greeter/greeter.proto.
