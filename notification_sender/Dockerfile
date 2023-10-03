# Start with a base image containing the Go runtime
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /APP

COPY . .
WORKDIR  /APP/notification_sender
RUN go build -o sender_binary
RUN mv sender_binary /
WORKDIR /
RUN rm -rf /APP
# Set the entry point for the container
CMD ["./sender_binary"]