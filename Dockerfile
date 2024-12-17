FROM golang:latest AS builder
RUN apt-get update && apt-get install -y make
WORKDIR /app
COPY . .
RUN make build
CMD ["make", "build-run"]
