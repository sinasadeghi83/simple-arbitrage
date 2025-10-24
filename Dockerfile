FROM golang:1.25.3-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/scraper

FROM scratch

WORKDIR /app

ENV METRICS_PORT=2112

EXPOSE ${METRICS_PORT}

COPY --from=builder /app /app

ENTRYPOINT [ "/app/scraper" ]
