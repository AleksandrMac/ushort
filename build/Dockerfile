FROM golang:1.17 as modules

RUN ls

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /src
ADD . /src
WORKDIR /src

# Добавляем непривилегированного пользователя
RUN useradd -u 10001 ushort

# Собираем бинарный файл
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
   go build -o /ushort ./cmd/ushort

FROM scratch

# Не забываем скопировать /etc/passwd с предыдущего стейджа
COPY --from=builder /etc/passwd /etc/passwd
USER ushort

COPY --from=builder /ushort /ushort
# COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["/ushort"]

# запуск docker build -f build/Dockerfile --tag aleksandrmac/ushort:latest .