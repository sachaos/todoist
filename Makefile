.PHONY: install
install: filter_parser.go
	go install

filter_parser.go: filter_parser.y
	goyacc -o filter_parser.go filter_parser.y

.PHONY: clean
clean:
	rm -f y.output

.PHONY: test
test:
	go test -v
