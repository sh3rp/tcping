all: dep install

dep:
	dep ensure

compile:
	go build cmd/tcping/tcping.go

install:
	go install cmd/tcping/tcping.go
	sudo setcap cap_net_raw+ep $(GOPATH)/bin/tcping

.PHONY: build