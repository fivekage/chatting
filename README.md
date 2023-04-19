# Stay Chatting

> Real time chatting API for the Stay webapp, using Websockets and made with Golang, sweat and tears.

## How to use

Well, for the moment, you can run it with `go run main.go`, browse to `localhost:${PORT}` and enjoy the ride.

## Docker

### Build :

`docker build -t stay.chatting .`

### Run :

`docker run -p 5000:8080 --env-file .env --name stay.chatting stay.chatting`
