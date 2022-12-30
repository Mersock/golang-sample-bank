#build
FROM golang:1.18.3-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#runner
FROM alpine:3.16 as runner
WORKDIR /app
#build app
COPY --from=builder /app/main .
#migrate db
COPY db/migration ./db/migration
#env
COPY app.env .
#start sh
COPY start.sh .
#copy wait for to synchronize docker container running
#ref https://github.com/mrako/wait-for
COPY wait-for.sh .

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]