.PHONY = build server

SHORTINHO = github.com/pinheirolucas/shortinho
BIN = ./bin/shortinho

build:
	go build -o ${BIN} ${SHORTINHO}

run:
	go run ${SHORTINHO}

server: build
	GIN_MODE=release ${BIN}

.PHONY = clean
clean:
	go clean
	rm -rf ./bin
