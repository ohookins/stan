.PHONY: build
build: Dockerfile
	docker build . -t ohookins/stan:latest

.PHONY: test
test: Dockerfile.test
	docker build -t ohookins/stan:test -f Dockerfile.test .
	docker run -it --rm ohookins/stan:test

.PHONY: run
run: build
	docker run -it --rm -p 8080:8080 ohookins/stan:latest