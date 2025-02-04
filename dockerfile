# Use the official Go image as the base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on
EXPOSE 8080

# Set the environment variable for the GitHub token
# (You should set the actual token value when running the container)
ENV GITHUB_TOKEN=""

# Command to run the application
CMD ["sh", "-c", "GITHUB_TOKEN=$GITHUB_TOKEN ./main"]
