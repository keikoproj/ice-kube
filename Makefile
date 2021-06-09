
# Image URL to use all building/pushing image targets
IMAGE_TAG ?= ice-kube:latest


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


all: test build

setup: ; $(info $(M) setting up env variables for test…) @ ## Setup env variables
export LOCAL=true

# Run tests
test: setup fmt vet
	go test ./... -coverprofile cover.out


.PHONY:build
build: fmt vet
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -tags release $(GO_LDFLAGS) -o bin/ice-kube main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Deploy ice-kube in the configured Kubernetes cluster in ~/.kube/config
deploy:
	kubectl apply -f Icekube.yaml

undeploy:
	kubectl delete -f Icekube.yaml


# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

.PHONY: lint
lint:
	go get -u golang.org/x/lint/golint
	@echo "golint $(LINTARGS)"
	@for pkg in $(shell go list ./...) ; do \
		golint $(LINTARGS) $$pkg ; \
	done

# Build the docker image
docker-build:
	docker build . -t ${IMAGE_TAG}

# Push the docker image
docker-push:
	docker push ${IMAGE_TAG}
