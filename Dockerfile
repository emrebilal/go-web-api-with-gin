FROM golang:1.19-bullseye

WORKDIR /source
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN swag init
RUN go build

EXPOSE 8080
CMD ["./rating-api"]