# Use the official Golang image from the Docker Hub
FROM golang:latest

# Create a directory inside the container to store all our application and then make it the working directory.
WORKDIR /app/bin/mdm

# Copy the service account key into the Docker image.
COPY my-key.json /app/bin/mdm/

# Set the GOOGLE_APPLICATION_CREDENTIALS environment variable.
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/bin/mdm/my-key.json

# Copy the entire project
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Expose port 8080 in the Docker image
EXPOSE 8080

# Run the command
CMD ["go", "run", "./cmd/micromdm/serve.go"]
