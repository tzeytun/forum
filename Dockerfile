# Building step
FROM golang:latest AS build

# Install dependencies
RUN apt-get update && apt-get install -y gcc libc6-dev

# Set working directory
WORKDIR /go/src/forum-project

# Copy the files into the container
COPY . .

# Install dependencies and compile the application
RUN go mod download && go mod verify && \
    go build -o forum-project .

# Running step
FROM debian:latest

# Install dependencies
RUN apt-get update && apt-get install -y ca-certificates

# Set working directory
WORKDIR /root/

# Copy the compiled application from the build stage
COPY --from=build /go/src/forum-project/forum-project .

# Copy frontend files
COPY --from=build /go/src/forum-project/frontend /root/frontend

# Identify the executable file
ENTRYPOINT ["./forum-project"]

# App runs on port 8080
EXPOSE 8080
