# https://docs.docker.com/engine/reference/builder/

# Use an official Python runtime as a parent image
# See https://hub.docker.com/ and search for python:2.7-slim
FROM golang:1.19

# Set the working directory to /app
# The WORKDIR instruction sets the working directory (inside the container) 
# for any RUN, CMD, ENTRYPOINT, COPY and ADD instructions that 
# follow it in the Dockerfile. 
WORKDIR /golang-url-shortner

# Copy the application directory (app) contents into the container at /app
ADD url-shortner /golang-url-shortner/

# Install any needed packages specified in requirements.txt
# RUN during image build
RUN go get
RUN go build

RUN mkdir -p logs
# Make port 80 available to the world outside this container
EXPOSE 9000

# Run python app.py when the container launches
# This happens if no command is specified
# CMD ["./url-shortner"]
CMD ["./url-shortner"]
# CMD ["tail", "-f", "/dev/null"]
# CMD ["python", "app.py"]
