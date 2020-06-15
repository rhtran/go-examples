# build executable binary
FROM golang:1.14.1-alpine as builder

COPY . $GOPATH/src/go-host-ip
WORKDIR $GOPATH/src/go-host-ip

# git
RUN apk add --no-cache git

# Install library dependencies
RUN go mod tidy

# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/go-host-ip

# build a small image
FROM scratch

# copy static executable
COPY --from=builder /go/bin/go-host-ip  /go/bin/go-host-ip

ENV PORT 8090
EXPOSE 8090

ENTRYPOINT ["/go/bin/go-host-ip"]