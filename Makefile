.PHONY: install
install: filter_parser.go
	go install

.PHONY: test
test: filter_parser.go
	go test -v

filter_parser.go: filter_parser.y
	goyacc -o filter_parser.go filter_parser.y

.PHONY: clean
clean:
	rm -f y.output filter_parser.go
