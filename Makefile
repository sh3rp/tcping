all: dep protos install

dep:
	dep ensure

protos:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I=rpc --go_out=plugins=grpc:rpc rpc/tcping.proto 

install:
	go install cmd/tcping/tcping.go
	go install cmd/tcpingd/tcpingd.go
	go install cmd/tpctl/tpctl.go
	sudo setcap cap_net_raw+ep $(GOPATH)/bin/tcping

.PHONY: install protos dep