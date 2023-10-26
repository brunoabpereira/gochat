FROM golang:1.20.4-alpine as build

WORKDIR /src

# install go dependencies
COPY go.mod go.sum .
RUN go mod download && go mod verify

# compile
COPY . .
RUN go build -v -o run 

###

FROM alpine:3.18

COPY --from=build /src/run /app/run

EXPOSE 9000

CMD ["/app/run"]