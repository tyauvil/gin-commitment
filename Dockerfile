# Build container
FROM golang:1.10-alpine as BUILD

ENV CGO_ENABLED=0
ARG VERSION=0.0.4
ARG SOURCE_COMMIT=""
ARG SOURCE_BRANCH=""

WORKDIR /go/src/app
COPY . .

RUN env
RUN apk add git --no-cache
RUN go get -d -v ./...
RUN go install -ldflags="-d -s -w -X main.SourceBranch=$SOURCE_BRANCH -X main.GolangVersion=$GOLANG_VERSION -X main.SourceCommit=$SOURCE_COMMIT" -v ./...

# Release container
FROM scratch as RELEASE

ENV GIN_MODE=release
ENV PORT=8080

COPY --from=BUILD /go/bin/* /
COPY --from=BUILD /etc/ssl/certs/ /etc/ssl/certs/
COPY ./static /static

EXPOSE 80 443 8080
ENTRYPOINT ["/app"]
