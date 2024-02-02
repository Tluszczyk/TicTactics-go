# Start from the official golang image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Create a directory for the build
RUN mkdir /app/src

# Copy the API build from the current directory to the working directory
COPY ../../../ /app/src/

RUN cd /app/src/IOService/APIGateway
RUN go build 

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD pwd && tail -f /dev/null
