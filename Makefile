.PHONY: all
all: process-env

.PHONY: clean
clean:
	-rm process-env

process-env: main.go
	CGO_ENABLED=0 go build -trimpath
