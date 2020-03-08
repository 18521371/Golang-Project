# Use the official image as a parent image
FROM ubuntu:18.04

# Run the command inside your image filesystem
RUN apt update && apt upgrade
RUN apt -y install curl && curl -V

#Download the archive and extract it into /usr/local
RUN curl -sSL https://storage.googleapis.com/golang/go1.14.linux-amd64.tar.gz \
		| tar -v -C /usr/local -xz
# Install nano to edit code easily inside container
RUN apt-get install nano 

# Create necessary folder 
RUN mkdir -p usr/src/app && mkdir -p usr/src/app/log

# Set up environment variable PATH
ENV PATH /go/bin:/usr/src/app:/usr/local/go/bin:$PATH

# Decide where to work at start
WORKDIR /usr/src/app

# Copy the file from your host to your current location
COPY Server.go .

# Execute Server.go in starting container 
CMD ["go","run","Server.go"]

# Inform Docker that the container is listening on the specified port at runtime.
EXPOSE 8080

