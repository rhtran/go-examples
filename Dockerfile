# build executable binary
FROM golang:1.14.1-alpine as builder

COPY . $GOPATH/src/go-examples
WORKDIR $GOPATH/src/go-examples

# git
RUN apk add --no-cache git

# Install library dependencies
RUN go mod tidy

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/go-examples

# build a small image
FROM scratch

# copy static executable
COPY --from=builder /go/bin/go-examples  /go/bin/go-examples

ENV PORT 8090
EXPOSE 8090

ENTRYPOINT ["/go/bin/go-examples"]