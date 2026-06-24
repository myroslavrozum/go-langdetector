.PHONY: clean clean-data clean-all 

templ-generate:
	go tool templ generate

test:
	go test ./...

build: templ-generate test
	@echo "Building the executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64
	go tool templ generate
	go build -o ./bin/go-langdetector .

build-static: test
	@echo "Building a static executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64
	go build -ldflags="-w -s" -o ./bin/go-langdetector-static .
	strip ./bin/go-langdetector-static

compress: build-static
	@echo "Compressing the executable...."
	upx --brute --force-macos bin/go-langdetector-static

run: build
	@echo "Running the executable...."
	./bin/go-langdetector

clean-bin:
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

clean-templ-generated:
	rm webapp/*_templ.go

clean: clean-bin clean-data clean-templ-generated
	go clean -modcache

docker: docker-build docker-run
