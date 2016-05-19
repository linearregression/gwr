PACKAGES=$(shell glide novendor)

.PHONY: lint

lint:
	go vet $(PACKAGES)

.PHONY: test

test: check-license lint
	find . -type f -name '*.go' -not -name '*_string.go' | xargs golint
	go test $(PACKAGES)

vendor: glide.lock
	glide install

vendor/github.com/uber/uber-licence/package.json: vendor

vendor/github.com/uber/uber-licence/node_modules: vendor/github.com/uber/uber-licence/package.json
	cd vendor/github.com/uber/uber-licence && npm install

.PHONY: check-license add-license

check-license: vendor/github.com/uber/uber-licence/node_modules
	./vendor/github.com/uber/uber-licence/bin/licence --dry --file '*.go'

add-license: vendor/github.com/uber/uber-licence/node_modules
	./vendor/github.com/uber/uber-licence/bin/licence --verbose --file '*.go'
