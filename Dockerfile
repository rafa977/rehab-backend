# Start with a base image containing the Go runtime
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /APP

COPY . .
WORKDIR  /APP
RUN go build -o main
RUN mv main /
WORKDIR /
RUN rm -rf /APP
# Set the entry point for the container
CMD ["./main"]