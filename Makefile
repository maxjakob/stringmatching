GOVER=go1.4.2

DEPS=$(shell pwd)/.deps
OS=$(shell uname)
ARCH=$(shell uname -m)
GOOS=$(subst Darwin,darwin,$(subst Linux,linux,$(OS)))
GOARCH=$(subst x86_64,amd64,$(ARCH))
GOPKG=$(subst darwin-amd64,darwin-amd64-osx10.8,$(GOVER).$(GOOS)-$(GOARCH).tar.gz)

GOPATH=$(DEPS)
GOROOT=$(DEPS)/go
GOENV=GOROOT=$(GOROOT) GOPATH=$(GOPATH) GOBIN=$(GOROOT)/bin
GO=$(GOENV) $(GOROOT)/bin/go

.PHONY: build
build: test
	$(GO) build -v -o stringmatching

.PHONY: test
test: dependencies
	$(GO) test -v

.PHONY: clean
clean:
	rm -rf $(DEPS) stringmatching

dependencies: $(DEPS)/$(GOPKG)
	$(GO) get -v

$(DEPS)/$(GOPKG):
	mkdir -p $(dir $@)
	curl -s -o $@ https://storage.googleapis.com/golang/$(GOPKG)
	tar -C $(DEPS) -xzf $@

