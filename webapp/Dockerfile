FROM golang:1.20.4-alpine as build

WORKDIR /src

# install npm
RUN apk add npm

# install go dependencies
COPY go.mod go.sum .
RUN go mod download && go mod verify

# install npm dependencies
COPY chat-client/package.json ./chat-client/package.json
RUN cd chat-client/ && npm install --force

# compile
COPY app.go  .
COPY chat-client ./chat-client
RUN go build -v -o run && \
    cd chat-client/ && npm run build

###

FROM alpine:3.18

ENV HTML_DIR="/app/build"

COPY --from=build /src/run /app/run

COPY --from=build /src/chat-client/build /app/build

EXPOSE 8000

CMD ["/app/run"]