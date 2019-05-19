name    := wptool
ver     := $(shell git describe --always --dirty)
app     := $(shell which $(name))
gofiles := $(wildcard *.go) $(wildcard cmd/$(name)/*.go) $(wildcard wpfeed/*.go) \
	$(wildcard wphtml/*.go)
ldflags := -ldflags "-X main.version=$(ver) -w -s"
binlist := wptool feed-utils recover

install: $(app)

release: $(app)
	tar -czf $(name)-$(ver).tgz -C $(GOBIN) $(binlist)

$(app): $(gofiles)
	go install $(ldflags) ./...

$(name): $(gofiles)
	cd cmd/$(name); go build $(ldflags) -o $(name) ./...
	mv cmd/$(name)/$(name) .

$(name)-mac: $(gofiles)
	cd cmd/$(name); GOOS=darwin go build $(ldflags) -o $(name)-mac ./...
	mv cmd/$(name)/$(name)-mac .
