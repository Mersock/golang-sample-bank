#build
FROM golang:1.18.3-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
#install migreate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz


#runner
FROM alpine:3.16 as runner
WORKDIR /app
#build app
COPY --from=builder /app/main .
#migrate db
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
#env
COPY app.env .
#start sh
COPY start.sh .
#copy wait for to synchronize docker container running
COPY wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]