build:
	cd src; go fmt .; go build -o go-service

linux:
	cd src; go fmt .; GOOS=linux go build -o go-service

run:
	./run.sh

test:
	./src/test/curl-tests.sh
