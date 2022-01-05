FROM debian:latest

RUN apt update

RUN apt install -y dirmngr gnupg apt-transport-https software-properties-common ca-certificates curl

RUN curl -fsSL https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -
RUN add-apt-repository 'deb https://repo.mongodb.org/apt/debian buster/mongodb-org/4.2 main'
RUN apt update
RUN apt install -y mongodb-org

# Define mountable directories.
VOLUME ["/data/db"]

# Define working directory.
WORKDIR /data

# Define default command.
CMD ["mongod"]

# Expose ports.
#   - 27017: process
#   - 28017: http
#   - 10000: messaging

EXPOSE 27017
EXPOSE 28017
EXPOSE 10000

COPY main /data/main

ENTRYPOINT ["/bin/bash", "-l", "-c", "/data/main"]