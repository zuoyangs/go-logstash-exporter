FROM harbor.mcdchina.net/infra/golang:1.24.5-alpine3.22 AS build

WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn \
    GOROOT=/usr/local/go/current \
    GOPATH=/root/ \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin:/usr/local/bin

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-logstash-exporter .

FROM harbor.mcdchina.net/infra/alpine:3.22 AS release

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=build /app/config.yaml ./config.yaml
COPY --from=build /app/go-logstash-exporter ./go-logstash-exporter
RUN chmod +x ./go-logstash-exporter

EXPOSE 8080
ENV EXPORTER_CONFIG_LOCATION=/app/config.yaml

RUN addgroup -S exporter && adduser -S exporter -G exporter
USER exporter

ENTRYPOINT ["./go-logstash-exporter"]
