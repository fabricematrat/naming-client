# Makefile for the naming-client

ifndef GOPATH
$(warning You need to set up a GOPATH.)
endif

# Check runs tests.
check:
	go test ./...

# Update the project Go dependencies to the required revision.
deps: $(GOPATH)/bin/godeps
	$(GOPATH)/bin/godeps -u dependencies.tsv

# Generate the dependencies file.
create-deps: $(GOPATH)/bin/godeps
	godeps -t $(shell go list ./...) > dependencies.tsv || true

# Clean cleans.
clean:
	go clean ./...

# Reformat source files.
format:
	gofmt -w -l .

# Reformat and simplify source files.
simplify:
	gofmt -w -l -s .

