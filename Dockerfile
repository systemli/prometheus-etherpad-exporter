FROM alpine:3.13.3 as builder

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


FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY prometheus-etherpad-exporter /prometheus-etherpad-exporter

USER appuser:appuser

EXPOSE 9011

ENTRYPOINT ["/prometheus-etherpad-exporter"]
