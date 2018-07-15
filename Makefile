all: dep install

dep:
	dep ensure

compile:
	go build cmd/tcping/tcping.go

protos:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I=rpc --go_out=grpc:rpc rpc/tcping.proto 

install:
	go install cmd/tcping/tcping.go
	sudo setcap cap_net_raw+ep $(GOPATH)/bin/tcping

.PHONY: build