FROM golang:1.16.2-alpine as builder

WORKDIR /go/src/github.com/systemli/prometheus-etherpad-exporter

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

ADD . /go/src/github.com/systemli/prometheus-etherpad-exporter
RUN go get -d -v && \
    go mod download && \
    go mod verify && \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o /prometheus-etherpad-exporter


FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /prometheus-etherpad-exporter /prometheus-etherpad-exporter

USER appuser:appuser

EXPOSE 9011

ENTRYPOINT ["/prometheus-etherpad-exporter"]
