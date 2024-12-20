FROM ubuntu:latest

RUN apt-get update && apt-get upgrade
ADD src/go-service /usr/local/bin

WORKDIR /usr/local/bin
EXPOSE 8081
CMD ["./go-service"]
