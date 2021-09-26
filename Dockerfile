# Start by building the application.
FROM golang:1.16-buster as build

WORKDIR /go/src/app
ADD src /go/src/app

RUN go get -d -v ./...

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /go/bin/app

# Now copy it into our base image.
FROM alpine
COPY --from=build /go/bin/app /
CMD ["sh", "-c", "/app"]