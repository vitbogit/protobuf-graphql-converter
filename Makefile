.PHONY: command clean

GRAPHQL_CMD=protoc-gen-graphql-schema
VERSION=$(or ${tag}, dev)
UNAME:=$(shell uname)

# TODO: Somehow deal with PROTOPATH.
# In fact PROTOPATH is not needed for distibuting or running the plugin,
# but PROTOPATH is needed for building the "core" of it, 
# which is graphql/graphql.pb.go file. So PROTOPATH is used 
# only in "make build" command". How should we parse PROTOPATH?
ifeq ($(UNAME), Darwin)
	PROTOPATH := $(shell brew --prefix protobuf)/include
endif
ifeq ($(UNAME), Linux)
	PROTOPATH := /usr/local/include
endif
PROTOPATH := C:\Users\vitbogit\Documents\protoc\include

# Distributes generator (creates executable file "protoc-gen-graphql-schema" in dist folder)
distribute: clean
	cd ${GRAPHQL_CMD} && \
		go build \
			-ldflags "-X main.version=${VERSION}" \
			-o ../dist/${GRAPHQL_CMD}

# It is temporary command for speeding up development.
# Runs example "Greeter", but first redistributes generator to make sure it`s up to date
example: distribute 
	make -C ./examples/greeter generate

# Creates graphql/graphql.pb.go file from include/graphql/graphql.proto schema
plugin:
	protoc -I ${PROTOPATH} \
		-I include/graphql \
		--go_out=./graphql \
		include/graphql/graphql.proto
	mv graphql/github.com/vitbogit/protobuf-graphql-converter/graphql/graphql.pb.go graphql/
	rm -rf graphql/github.com

# TODO: Linter (old linter didn`t work)

# Testing
test:
	go list ./... | xargs go test

# TODO: Build (old logic for testing code in build command didn`t work well)

# Cleans executables of generator
clean:
	rm -rf ./dist/*

# TODO: Make sure it works as expected
all: clean build
	cd ${GRAPHQL_CMD} && GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ../dist/${GRAPHQL_CMD}.darwin
	cd ${GRAPHQL_CMD} && GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o ../dist/${GRAPHQL_CMD}.linux

