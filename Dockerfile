FROM golang:1.18-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o app.out

WORKDIR /dist
RUN cp /app/app.out .
RUN chmod +x /dist/app.out
RUN ls /dist

EXPOSE 4001

FROM scratch

COPY --from=builder /dist/app.out .
ENTRYPOINT [ "./app.out" ]

