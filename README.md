# Message-Rest-API

RESTful API that allows storage of messages and retrieval/modification through docker container. The following URLs can be used in order to access the functionality:

```
- {IP_ADDRESS}/messages (GET) (Will retrieve all messages)
- {IP_ADDRESS}/messages/{message_id} (GET) (Will retrieve the message with the specified ID)
- {IP_ADDRESS}/messages (POST) (Will create a new message)
- {IP_ADDRESS}/messages/{message_id} (PUT) (Include body for replacement/modification of message)
- {IP_ADDRESS}/messages/{mesage_id} (DELETE) (Allows for removal of specific message)
```

This program can run 1 of 2 ways - either with the embedded mongoDB or with a separate mongoDB docker container. 

## With embedded MongoDB:

``` docker run --net message-net --ip 172.18.0.2 --publish 10000:10000 --publish 27017:27017 --name messaging-api --restart always thavarti/messaging:latest ```

## With separate MongoDB:

- Firstly, a network should be created allowing for static IPs and the docker containers to communicate
- ``` docker network create --subnet=172.18.0.0/16 message-net ```
- _The subnet and name (message-net) can be substituted for any values of your choosing, however in this example message-net will be used._
- A MongoDB docker container should be ran on the network that has just been created, this can be done like so:
- ``` docker run --net message-net --ip 172.18.0.5 --name mongodb -d mongo:latest ```
- _Keep in mind the IP address set here will affect the environment variable in the next step_
- Next, run the docker command. keep in mind this time there is no publish for the 27017 port as we will be using the other container for our DB.
- ``` docker run --net message-net --ip 172.18.0.2 --env MONGODB=172.18.0.5 --publish 10000:10000 --name messaging-api --restart always messaging ```

In order to create a message, the following json format should be used:

```
{
    "content": "message",
    "author": {
        "firstname": "Taylor",
        "lastname": "Havart"
    }
}
```

_If an ID or Timestamp is entered, they will be overriden by the application_

## Testing

- Postman is the preferred tool to test this API calls, this can be done here: https://www.postman.com/downloads/
- A header will need to be added for KEY: Content-Type and VALUE: application/json for any POST requests
- Body tab can be used for json request data
- https://learning.postman.com/docs/getting-started/sending-the-first-request/

## Docker

- Container can be found here: https://hub.docker.com/repository/docker/thavarti/messaging