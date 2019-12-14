FROM golang:1.13 as serverbuilder
WORKDIR /server
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM scratch

COPY --from=serverbuilder /server/app  /server/app

WORKDIR /server

CMD ["./app"]
