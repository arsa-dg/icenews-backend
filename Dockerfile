FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /build

EXPOSE 8080

# Run the executable
CMD ["/build"]