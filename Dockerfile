FROM alpine:latest

RUN mkdir -p /usr/src/app

RUN apk update && apk upgrade
ADD src/go-service /usr/src/app

WORKDIR /usr/src/app
EXPOSE 8081
CMD ["./go-service"]
