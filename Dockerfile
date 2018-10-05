FROM golang:1.11 as builder

LABEL maintainer "Ryan Gurney <rygurney@zendesk.com>"

ENV GOPATH /go

WORKDIR /go/src/github.com/ragurney/samson-job

COPY . .

RUN GIT_REV=$(git rev-parse HEAD) \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && CGO_ENABLED=0 GOOS=linux \
	go build -a -installsuffix cgo \
	-o app -ldflags "-X main.GitRevision=${GIT_REV}" .

FROM debian:buster

LABEL maintainer "Ryan Gurney <rygurney@zendesk.com>"

WORKDIR /app

COPY --from=builder /etc/ssl /etc/ssl

COPY --from=builder /go/src/github.com/ragurney/samson-job/app .

ENTRYPOINT ["/app/app"]