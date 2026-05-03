cli-build:
	go build -o rune ./cmd/cli
cli-install: cli-build
	sudo mv rune /usr/local/bin/rune

server-build:
	go build -o rune-server ./cmd/server

server-install: server-build
	 sudo mv rune-server /usr/local/bin/rune-server