FROM golang:1.21.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o radish .
EXPOSE 6379
CMD ["./radish"]