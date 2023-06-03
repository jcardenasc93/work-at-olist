build-authors:
	@cd cmd/authors && go build -o ../../bin/load_authors

build:
	@cd app/ && go build -o ../bin/app

run: build
	@./bin/app

test:
	@go test -v ./... --cover

