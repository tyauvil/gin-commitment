# Build container
FROM golang:1.9 as BUILD

ENV CGO_ENABLED=0

WORKDIR /go/src/commit
COPY . .

RUN go get -d -v ./...
RUN go install -ldflags="-s -w" -v ./...

# Release container
FROM scratch as RELEASE

ENV GIN_MODE=release
ENV PORT=8080

COPY --from=BUILD /go/bin/* /
COPY ./commit_messages.txt ./names.txt /

ENTRYPOINT ["/commit"]
