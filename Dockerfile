# Use the official Golang image as the base image
FROM golang:1.21-alpine

WORKDIR /GoWeatherAPI

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download



COPY . .

RUN go build .

EXPOSE 8090

CMD ls
CMD ["./GoWeatherAPI serve"]

