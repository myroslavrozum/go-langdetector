build:
	@echo "Building the executable...."
	go build -o ./bin/go-langdetector .

build-static:
	@echo "Building a static executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64
	go build -ldflags="-w -s" -o ./bin/go-langdetector-static .

compress: build-static
	@echo "Compressing the executable...."
	upx --brute --force-macos bin/go-langdetector-static

run: build
	@echo "Running the executable...."
	./bin/go-langdetector

clean:
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

docker: docker-build docker-run
