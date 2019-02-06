FROM golang:1.10-alpine as builder
WORKDIR /go/src/bfx
COPY . .
RUN \ 
	apk add git && \
	go get -d -v ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags="-w -s" -v ./...

FROM scratch
COPY --from=builder /go/bin/bfx /
WORKDIR /
CMD ["/bfx"]
