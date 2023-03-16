all: build generate

build: ## compile binaries
	CGO_ENABLED=0 go build \
	-o bin/openstack-cni cmd/openstack-cni/main.go
	CGO_ENABLED=0 go build \
	-o bin/openstack-cni-daemon cmd/openstack-cni-daemon/main.go

generate: ## generate mocks
	go install github.com/matryer/moq@latest
	go generate ./...

clean: ## remove binaries
	go clean
	rm -f bin/*

.PHONY: test
test: ## run all tests
	go test -v -shuffle=on ./...
