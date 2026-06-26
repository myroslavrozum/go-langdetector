CURRENT_USER := $(shell whoami)
HEAD_HASH := $(shell git rev-parse HEAD)

.PHONY: makeversion clean clean-data clean-all 

build: makeversion generate test
	@echo "Building the executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64

	go build -o ./bin/go-langdetector .

makeversion:
	@echo $(HEAD_HASH)
	@echo $(HEAD_HASH) > .version

generate:
	go generate ./...

test:
	go test ./...

build-static: makeversion generate test
	@echo "Building a static executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64
	go build -ldflags="-w -s" -o ./bin/go-langdetector-static .
	strip ./bin/go-langdetector-static

compress: build-static
	@echo "Compressing the executable...."
	upx --brute --force-macos bin/go-langdetector-static

run:
	@echo "Running the app...."
	@echo "DEV" > .version
	go generate ./...
	go run .

clean-bin:
	@echo "DEV" > .version
	rm -f ./bin/*

all: build run clean

all-static: compress
	@echo "..... and running it"
	bin/go-langdetector-static

docker-build: docker-clean
	docker build -t go-langdetector:latest .

docker-run:
	docker run -d --rm -it \
		-u $(id -u):$(id -g) \
		-v $(CURDIR)/data:/app/data -p 8080:8080 \
		--name go-langdetector \
		go-langdetector:latest

docker-clean:
	docker rm -f go-langdetector
	docker rmi -f go-langdetector

clean-data:
	rm -rdf ./data/*

clean-generated:
	-rm webapp/*_templ.go
	-rm webapp/*_compiled.css

clean: clean-bin clean-data clean-generated
	go clean -modcache

docker: docker-build docker-run

deps:
	go get -u github.com/dgraph-io/badger/v4
	go get -u github.com/gin-gonic/gin
	go get -tool github.com/a-h/templ/cmd/templ@latest
	go get -u github.com/gorilla/websocket
	brew install aureuma/tailwindcss/tailwindcss-standalone
	go mod tidy