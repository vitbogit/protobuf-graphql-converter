.PHONY: build

build:
	protoc \
	  -I. \
		-I../../include/graphql \
		--plugin=../../dist/protobuf-to-graphql \
		--graphql_out=verbose:./greeter \
	  greeter.proto

# build:
# 	protoc \
# 	  -I. \
# 		-I../../include/graphql \
# 		--plugin=../../dist/protoc-gen-graphql \
# 		--graphql_out=verbose:./greeter \
# 	  greeter.proto