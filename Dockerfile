FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o crm.shopgo ./cmd/server

FROM scratch

COPY ./config /config

COPY --from=builder /build/crm.shopgo /

ENTRYPOINT [ "/crm.shopgo", "config/local.yaml" ]
