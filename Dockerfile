# Use the official Golang image as a base
FROM golang:1.23.4

# Set the working directory
WORKDIR /Task1

# Copy the Go module files and download dependencies
COPY go.mod ./ 


# Copy the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
