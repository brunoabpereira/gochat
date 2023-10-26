FROM golang:1.20.4-alpine

WORKDIR /src

RUN apk add npm

COPY go.mod go.sum .
RUN go mod download && go mod verify

COPY app.go .

COPY chat-client ./chat-client

RUN go build -v -o /app/run && \
    cd chat-client/ && npm run build && cp -r /src/chat-client/build /app/ && \
    rm -rf /src/*

EXPOSE 8000

ENV HTML_DIR="/app/build"

CMD ["/app/run","/app/build"]