build:
	go build -o rune ./cmd/cli
install: build
	sudo mv rune /usr/local/bin/rune