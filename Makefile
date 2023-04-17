all: build

build: agent

agent:
	env CGO_ENABLED=0 go build -o bin/client ./cmd/client/

clean:
	rm -f ./bin/client
