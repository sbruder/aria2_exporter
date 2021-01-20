FROM golang:alpine as builder

WORKDIR /go/src/github.com/sbruder/aria2_exporter/

COPY . .

RUN apk add --no-cache git

RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" .

FROM scratch

COPY --from=builder /go/src/github.com/sbruder/aria2_exporter/aria2_exporter /aria2_exporter

USER 1000

ENTRYPOINT ["/aria2_exporter"]

EXPOSE 9578
