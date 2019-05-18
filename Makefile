name    := wptool
ver     := $(shell git describe --always --dirty)
app     := $(shell which $(name))
gofiles := $(wildcard *.go) $(wildcard cmd/$(name)/*.go) $(wildcard wpfeed/*.go)
ldflags := -ldflags "-X main.version=$(ver) -w -s"

install: $(app)

$(app): $(gofiles)
	go install $(ldflags) ./...

$(name): $(gofiles)
	cd cmd/$(name); go build $(ldflags) -o $(name) ./...
	mv cmd/$(name)/$(name) .

$(name)-mac: $(gofiles)
	cd cmd/$(name); GOOS=darwin go build $(ldflags) -o $(name)-mac ./...
	mv cmd/$(name)/$(name)-mac .
