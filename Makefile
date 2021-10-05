GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto:
	time protoc \
		--proto_path=${GOPATH}/pkg/mod \
 		--proto_path=${GOPATH}/src \
 		--proto_path=. \
 		--micro_out=. \
 		--gogofaster_out=. proto/*.proto
