FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
#ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/zero

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY http/etc /app/etc
COPY data /data
RUN go build -ldflags="-s -w" -o /app/http http/http.go


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/http /app/http
COPY --from=builder /app/etc /app/etc
COPY --from=builder /data /data

EXPOSE 4001

CMD ["./http", "-f", "etc/http-api.yaml"]
