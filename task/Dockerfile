FROM golang:1.19.1-alpine3.15 AS builder 
COPY . /go/src/projects/LDmitryLD/task-service/task
WORKDIR /go/src/projects/LDmitryLD/task-service/task

RUN go build -ldflags="-w -w" -o /go/bin/server /go/src/projects/LDmitryLD/task-service/task/cmd/app

FROM alpine:3.15

COPY --from=builder /go/bin/server /go/bin/server 
COPY ./.env /app/.env

WORKDIR /app

ENTRYPOINT [ "/go/bin/server"]