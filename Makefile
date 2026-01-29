.PHONY: build clean install

build:
	go build -o n8r ./cmd/n8r

clean:
	rm -f n8r

install:
	go install ./cmd/n8r
