NAME = "gqlclientgen"

build:
	go build -o bin/$(NAME) ./cmd/$(NAME)

install:
	go install ./cmd/$(NAME)