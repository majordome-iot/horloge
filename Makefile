TARGET=./cmd/horloge/horloge

$(TARGET):
	cd ./cmd/horloge && go build .

build: clean $(TARGET)

run: build
	exec $(TARGET) run --sync redis

clean:
	rm -f $(TARGET)

docker-build:
	docker build . -t shinuza/horloge:latest

test:
	GIN_MODE=release go test

.PHONY: build