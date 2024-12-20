FROM golang:alpine

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# CMD ['go', 'run', 'cmd/api/main.go']