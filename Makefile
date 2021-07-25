rootdir = $(realpath .)
BUILD = go build -o
RUN = go run
.PHONY := all main run clean

all: main

main: $(rootdir)/src/main.go
	$(BUILD) $(rootdir)/bin/$@ $^

run: $(rootdir)/src/main.go
	$(RUN) $^

clean:
	rm -f  $(rootdir)/bin/main