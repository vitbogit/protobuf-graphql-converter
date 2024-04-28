.PHONY: generate

# Path to plugin`s executable
PLUGIN_DIST_PATH := ../../dist/protoc-gen-graphql-schema

# Generate graphql schema from protobuf schema
generate:
	protoc \
	  -I. \
		-I../../include/graphql \
		--plugin=$(PLUGIN_DIST_PATH) \
		--graphql-schema_out=./ \
	  greeter.proto
