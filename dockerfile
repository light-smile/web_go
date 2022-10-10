FROM alpine:latest

WORKDIR /app

COPY . .

RUN ./app_linux_armv7

EXPOSE 3000
