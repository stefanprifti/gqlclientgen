NAME = "gqlclientgen"

build:
	go build -o bin/$(NAME) ./cmd/$(NAME)

run:
	cd ./cmd/$(NAME)/testdata
	./bin/$(NAME)

install:
	go install ./cmd/$(NAME)