FROM alpine:3.11 as certs

RUN apk add --no-cache ca-certificates

FROM golang:1.13-alpine as builder

WORKDIR /go/src/app

ADD go.mod /go/src/app
ADD go.sum /go/src/app
RUN go mod download

ADD . /go/src/app
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION} -X main.date=$(date '+%FT%T.%N%:z')" -o /pushit

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /pushit /pushit

ARG VERSION=dev
LABEL version="${VERSION}" maintainer="fopina <https://github.com/fopina/pushit/>"

ENTRYPOINT [ "/pushit" ]
CMD [ "-w" ]
