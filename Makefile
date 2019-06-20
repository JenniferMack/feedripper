name    := feedpub
ver     := $(shell git describe --tags --always --dirty)
gofiles := *.go
ldflags := -ldflags "-X main.version=$(ver) -w -s"
relDir  := releases
.PHONY: install test

install:
	go install $(ldflags) ./...

test:
	go test $$(go list ./...)

pkg: $(gofiles) | $(relDir)
	cd cmd/$(name); go build -o ../../$(name) $(ldflags) ./...
	tar -czf $(relDir)/$(name)-$(ver).tgz $(name)
	rm $(name)

$(name): $(gofiles)
	cd cmd/$(name); go build $(ldflags) -o $(name) ./...
	mv cmd/$(name)/$(name) .

$(name)-mac: $(gofiles)
	cd cmd/$(name); GOOS=darwin go build $(ldflags) -o $(name)-mac ./...
	mv cmd/$(name)/$(name)-mac .

$(relDir): 
	mkdir $@
