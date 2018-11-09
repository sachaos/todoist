.PHONY: install
install: prepare
	go install

.PHONY: build
build: prepare
	go build

.PHONY: test
test: prepare
	go test -v

.PHONY: prepare
prepare: filter_parser.y
	go get golang.org/x/tools/cmd/goyacc
	goyacc -o filter_parser.go filter_parser.y

.PHONY: clean
clean:
	rm -f y.output filter_parser.go
