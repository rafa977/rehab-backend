# Start with a base image containing the Go runtime
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /APP

COPY . .
WORKDIR  /APP/mailer
RUN go build -o mailer_binary
RUN mv mailer_binary /
RUN mv .env /
RUN mv templates /
WORKDIR /
RUN rm -rf /APP
# Set the entry point for the container
CMD ["./mailer_binary"]