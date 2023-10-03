# Start with a base image containing the Go runtime
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /APP

COPY . .
WORKDIR  /APP/notification_receiver
RUN go build -o receiver_binary
RUN mv receiver_binary /
WORKDIR /
RUN rm -rf /APP
# Set the entry point for the container
CMD ["./receiver_binary"]