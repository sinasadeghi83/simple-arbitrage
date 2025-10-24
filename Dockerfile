FROM golang:1.25.3-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o /app/scraper

FROM scratch

ENV PORT=3000

EXPOSE ${PORT}

COPY --from=builder /app/scraper /app/scraper

ENTRYPOINT [ "/app/scraper" ]