build:
	cd src; go fmt .; go build -o go-service

run:
	cd src; ./go-service
