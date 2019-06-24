all: clean protos mod build install

clean:
	rm -rf include bin target readme.txt

mod:
	go mod tidy

protos:
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc -I=rpc --go_out=plugins=grpc:rpc rpc/tcping.proto 

build:
	mkdir target
	go build -o target/tcping cmd/tcping/tcping.go
	go build -o target/tcpingd cmd/tcpingd/tcpingd.go
	go build -o target/tpctl cmd/tpctl/tpctl.go
	cp target/tcping $(GOPATH)/bin/tcping

install:
	sudo setcap cap_net_raw+ep $(GOPATH)/bin/tcping

.PHONY: clean deps build install protos
