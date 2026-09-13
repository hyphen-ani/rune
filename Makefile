build-ui:
	cd ui && npm install && npm run build

cli-build:
	go build -o rune ./cmd/cli

cli-install: cli-build
	sudo mv rune /usr/local/bin/rune

server-build: build-ui
	go build -o rune-server ./cmd/server

server-install: server-build
	sudo mv rune-server /usr/local/bin/rune-server

make-server: build-ui
	go build -o rune-server ./cmd/server && \
	sudo mv rune-server /usr/local/bin/rune-server

make-cli:
	go build -o rune ./cmd/cli && \
	sudo mv rune /usr/local/bin/rune
