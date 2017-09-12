.PHONY: install
install: filter_parser.go
	go install

filter_parser.go:
	goyacc -o filter_parser.go filter_parser.y
