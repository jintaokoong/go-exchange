all: build

build:
	go build -o bin/exch main.go
clean:
	rm -r bin/
	rm *.log