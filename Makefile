BINARY_NAME=ptpr

.PHONY: build clean

all: build

build:
	go build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
