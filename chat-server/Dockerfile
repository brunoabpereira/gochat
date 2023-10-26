FROM golang:1.20.4-alpine

WORKDIR /src

COPY go.mod go.sum .
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /app/run && \
    rm -rf /src/*

EXPOSE 9000

CMD ["/app/run"]