.PHONY: generate

PLUGIN_DIST_PATH := ../../dist/protoc-gen-graphql-schema

generate:
	protoc \
	  -I. \
		-I../../include/graphql \
		--plugin=$(PLUGIN_DIST_PATH) \
		--graphql-schema_out=./ \
	  greeter.proto
