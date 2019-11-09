FROM golang:1.13.4 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build

FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./go-auto-yt"]
