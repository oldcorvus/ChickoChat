# Chicko Chat


Chicko chat is a real time multi room multi user chat using golang, websocket , mongodb with restful api

[Sample Front]

(https://github.com/younes-nb/chicko-frontend)



## Install ##
All that you need is [Golang](https://golang.org/).
```sh
 $ git clone https://github.com/oldcorvus/chicko_chat.git

```
```sh
$ go mod tidy

```
```sh
$ go build -o chat 

```
```sh
$ ./chat --mongoURI="mongodb://localhost:27017"

```
And navigate to `http://127.0.0.1:8080/`.


## Running Locally with Docker

1.build the image:

```sh
  $ docker-compose build .
```
2.Spin up the containers

```sh
  $ docker-compose up
```
And navigate to `http://127.0.0.1:8080/`.

## ChatRoom API ##

  * `GET /chats/`: to join a room and start chat 
  * `GET /ws/:roomId/:userId/`: to start a websocket connection

  * `POST /start/`: register and obtain user id 
  * `POST /user-rooms/`: ro get user rooms based on user id 
  * `POST /room-history/`: to retrieve messages of room based on room id 
  * `POST /add-user-room/`: to register user in room
  * `POST /room-user-details/`: to retrieve user data of a room

## Test ##

to run tests 

Controllers :
```sh
 $ cd controllers 

```
```sh
$ go test

```
Database 

```sh
 $ cd database 

```
```sh
$ go test

```
