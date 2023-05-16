# Use the official Golang image from the Docker Hub
FROM golang:latest

# Create a directory inside the container to store all our application and then make it the working directory.
WORKDIR /app/bin/mdm

# Copy the entire project
COPY . .


# Download all the dependencies
RUN go get -d -v ./...

# Run the command
CMD ["go", "run", "./cmd/micromdm/serve.go"]
