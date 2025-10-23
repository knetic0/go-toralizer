all: build

build:
	go build -o toralize toralize.go
clean:
	rm -f toralize
install: build
	mv toralize /usr/local/bin/toralize
uninstall:
	rm -f /usr/local/bin/toralize
.PHONY: all build clean install uninstall