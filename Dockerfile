FROM golang:1.12 as serverbuilder
COPY . .
WORKDIR /server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM scratch

COPY --from=serverbuilder /server/main  /server

WORKDIR  /server

EXPOSE 8091

CMD ["./main"]