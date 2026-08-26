cli-build:
	go build -o rune ./cmd/cli
cli-install: cli-build
	sudo mv rune /usr/local/bin/rune

server-build:
	go build -o rune-server ./cmd/server

server-install: server-build
	 sudo mv rune-server /usr/local/bin/rune-server

make-server:
	go build -o rune-server ./cmd/server && \
	sudo mv rune-server /usr/local/bin/rune-server

make-cli:
	go build -o rune ./cmd/cli && \
	sudo mv rune /usr/local/bin/rune

