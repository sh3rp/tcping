all: protos build install

protos:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I=rpc --go_out=plugins=grpc:rpc rpc/tcping.proto 

build:
	go build -o $(HOME)/bin/tcping cmd/tcping/tcping.go
	go build -o $(HOME)/bin/tcpingd cmd/tcpingd/tcpingd.go
	go build -o $(HOME)/bin/tpctl cmd/tpctl/tpctl.go

install:
	sudo setcap cap_net_raw+ep $(HOME)/bin/tcping

.PHONY: build install protos