FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src

COPY ["go.mod","go.sum","./"]

RUN go mod download

COPY . ./

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

COPY .env /.env

COPY files /files

EXPOSE 3333

CMD [ "/app" ]

