build-server:
	go build -o ./bin/server ./cmd/server

build-worker:
	go build -o ./bin/worker ./cmd/worker