FROM golang:1.16
WORKDIR /app
RUN mkdir -p /app/logs
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o openseasync main/main.go
EXPOSE 8080
VOLUME /app/logs
CMD ["./openseasync"]