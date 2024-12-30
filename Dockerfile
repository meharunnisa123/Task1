# Stage 1: Build
FROM golang:1.22.7 as builder
WORKDIR /app
COPY go.mod ./ 
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o main .

# Stage 2: Run
FROM frolvlad/alpine-glibc:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
