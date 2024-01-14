FROM golang:1.21.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mapdb .
ENV PORT 6379
ENV SERVER_TYPE 0
CMD ./mapdb --port $PORT --server-type $SERVER_TYPE