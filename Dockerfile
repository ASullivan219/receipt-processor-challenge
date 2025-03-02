FROM golang:1.23 as builder


# Copy Files to docker container /app
# Download all dependencies
WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

# Work in the /cmd directory
# build main.go
# Run the project exposing port 8080
WORKDIR /app/cmd
RUN go build -o ./cmd
EXPOSE 8080

WORKDIR /app
CMD ["./cmd/cmd"]
