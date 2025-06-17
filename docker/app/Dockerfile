FROM golang:1.24

RUN go install github.com/cespare/reflex@v0.3.1 \
    && go install github.com/go-delve/delve/cmd/dlv@v1.24.2 

COPY reflex.conf /
COPY start.sh /

ENTRYPOINT ["reflex", "-c", "/reflex.conf"]