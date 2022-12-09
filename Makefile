# Copyright 2022 Ciena Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.DEFAULT_GOAL = help

.PHONY: help
help: ## Display this message.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

VERSION=$(shell head -1 VERSION)
COMMIT=$(shell git log -n 1 --pretty="%H") $(shell git diff --quiet || echo "(dirty)")

build: ## Build the gns3ctl command
	@echo "VERSION: $(VERSION)"
	@echo "COMMIT: $(COMMIT)"
	CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ciena/gns3ctl/cmd.Version=$(VERSION) \
		          -X 'github.com/ciena/gns3ctl/cmd.Commit=$(COMMIT)'" \
		-o ./gns3ctl ./main.go

build-releases: ## Build multiple architecture releases
	@echo "VERSION: $(VERSION)"
	@echo "COMMIT: $(COMMIT)"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ciena/gns3ctl/cmd.Version=$(VERSION) \
		          -X 'github.com/ciena/gns3ctl/cmd.Commit=$(COMMIT)'" \
		-o ./gns3ctl-linux-amd64-$(CI_COMMIT_TAG) ./main.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ciena/gns3ctl/cmd.Version=$(VERSION) \
		          -X 'github.com/ciena/gns3ctl/cmd.Commit=$(COMMIT)'" \
		-o ./gns3ctl-linux-arm64-$(CI_COMMIT_TAG) ./main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ciena/gns3ctl/cmd.Version=$(VERSION) \
		          -X 'github.com/ciena/gns3ctl/cmd.Commit=$(COMMIT)'" \
		-o ./gns3ctl-darwin-amd64-$(CI_COMMIT_TAG) ./main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ciena/gns3ctl/cmd.Version=$(VERSION) \
		          -X 'github.com/ciena/gns3ctl/cmd.Commit=$(COMMIT)'" \
		-o ./gns3ctl-darwin-arm64-$(CI_COMMIT_TAG) ./main.go

clean: ## Remove any build artifacts
	rm -rf ./gns3ctl ./gns3ctl-*

format:
	CGO_ENABLED=0 go fmt ./...

vet:
	CGO_ENABLED=0 go vet ./...

test: ## Run unit tests
	CGO_ENABLED=0 go test ./...

.PHONY: lint
lint: ## Perform basic lint
	golangci-lint run -v --config=./.golangci.yaml

.PHONY: docker-lint
docker-lint: ## Perform basic lint via Docker
	docker run -i --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v --config=./.golangci.yaml

.PHONY: lint-pedantic
lint-pedantic: ## Perform pedantic lint
	golangci-lint run --config=/app/.golangci-pedantic.yaml

.PHONY: docker-lint-pedantic
docker-lint-pedantic: ## Perform pedantic lint
	docker run -ti --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run --config=/app/.golangci-pedantic.yaml
