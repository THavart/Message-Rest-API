# TAK Server Manager
_This project was built as a layer over docker containers in order to provide an easy way for ATAK users to create, manage, and interact with TAK Servers_

## Running and Testing

- This project was built in go, which requires a development environment to be created.
    - Visual studio code is the preferred IDE due to it's simplicity
    - Go must also be installed on your platform of choice: https://golang.org/doc/install
    - Theoretically this can run on any platform, however has only been tested on linux.
    - In order to interact with the docker containers, docker must also be installed: https://docs.docker.com/engine/install/ubuntu/
- To create an executable, run the command from the git repo: `go build -o build/tak-server-manager`

## Functionality

- Using a RESTful API architecture the program is designed to expose a http server that can accept commands from other clients, as well as provide useful webpages for connection. The following are supported by this software:
    - GET -> http://192.168.2.89:10000/servers/{id}
        - Using the ID of the server will retrieve info on that server
    - POST -> http://192.168.2.89:10000/servers
        - ```json
           // Note how for the code block the "id" and "mongoid" are left empty. There is also an absent "NetworkInfo" field.
            This is due to autopopulation upon creation.
            {
            "id": "",
            "mongoid": "",
            "name": "TAK-SERVER-5",
            "desc": "testing 1234",
            "author": {
                "firstname": "Jeff",
                "lastname": "Bezos",
                "dept": "RCMP"
                }
            }
          ```
    - DELETE ->  http://192.168.2.89:10000/servers/{id}
        - Deletion of servers will stop running containers, remove containers, and delete from the database.
    - Return values:
        - ```json
            {
                "id": "97962430269",
                "mongoid": "",
                "name": "TAK-SERVER-6",
                "desc": "testing 1234",
                "port": {
                    "exposedport": 10010,
                    "map_8080": 10006,
                    "map_8087": 10007,
                    "map_8443": 10009,
                    "map_8089": 10008,
                    "IPAddress": "172.20.0.7"
                },
                "author": {
                    "firstname": "Jeff",
                    "lastname": "Bezos",
                    "dept": "RCMP"
                }
            }
          ```
    - External webpages:
        - /guide/{key} provides detailed instructions on how to integrate a tak-server into ATAK
        - /download/{key} provides a .zip file of a tak-server data package

## Testing

- Postman is the preferred tool to test this API calls, this can be done here: https://www.postman.com/downloads/
- A header will need to be added for KEY: Content-Type and VALUE: application/json for any POST requests
- Body tab can be used for json request data
- https://learning.postman.com/docs/getting-started/sending-the-first-request/
