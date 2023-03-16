IMAGE_REPO?=MISSING_IMAGE_REPO
IMAGE_VERSION?=0.0.1
IMAGE_NAME?=openstack-cni

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

docker-build: ## build the binary in a container
	docker build -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION) .

docker-push: ## push the container into a repo
	docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_VERSION)

docker-release: docker-build docker-push ## build and push the container

helm-release:
	helm upgrade openstack-cni helm --install