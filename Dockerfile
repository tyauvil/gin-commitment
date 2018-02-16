# Build container
FROM golang:1-alpine as BUILD

ENV CGO_ENABLED=0

WORKDIR /go/src/commit
COPY . .

RUN apk add git --no-cache
RUN go get -d -v ./...
RUN go install -ldflags="-d -s -w" -v ./...

# Release container
FROM scratch as RELEASE

ENV GIN_MODE=release
ENV PORT=8080

COPY --from=BUILD /go/bin/* /
COPY ./commit_messages.txt ./names.txt ./index.tmpl /

ENTRYPOINT ["/commit"]
