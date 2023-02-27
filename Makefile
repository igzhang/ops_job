all: build

build: server agent

agent:
	env CGO_ENABLED=0 go build -o bin/client ./cmd/client/

server:
	env CGO_ENABLED=0 go build -o bin/server ./cmd/server/

clean:
	rm -f ./bin/client
	rm -f ./bin/server
