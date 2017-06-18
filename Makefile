.PHONY: all clean images

GOFILES := $(shell find . -type f -name '*.go')

all: bin/message-bot

bin/message-bot: $(GOFILES)
	mkdir -p bin
	cd cmd/message-bot && \
	CGO_ENABLED=0 GOOS=linux go build \
		-a -ldflags '-w' \
		-o ../../bin/message-bot .

images: bin/message-bot
	docker build -t message-bot:latest -f Dockerfile-message-bot .

clean:
	rm -fr bin
