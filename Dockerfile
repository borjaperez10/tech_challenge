  
FROM golang:1.12-alpine
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/tech_challenge

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/server server.go

# Run the binary program 
CMD ["./out/server"]