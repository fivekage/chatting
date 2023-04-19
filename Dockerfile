# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /docker-app

FROM alpine:latest as production 
COPY --from=build /docker-app /docker-app

ARG API_BASE_URL=http://localhost:5001

ENV PORT=8080
ENV API_BASE_URL=${API_BASE_URL}

EXPOSE 8080

CMD [ "/docker-app" ]