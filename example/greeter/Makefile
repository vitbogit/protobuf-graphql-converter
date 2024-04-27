.PHONY: generate

generate:
	protoc \
	  -I. \
		-I../../include/graphql \
		--plugin=../../dist/protoc-gen-graphql-schema \
		--graphql-schema_out=./ \
	  greeter.proto
