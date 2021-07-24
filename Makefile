.PHONY := all main run clean

all: main

main: main.go
	go build -o bin/$@ $^

run: main.go
	go run $^

clean:
	rm -f bin/main