# Chatbox with end-to-end encryption

## Program Execution

* Create a directory in Home
```
$ mkdir go
```
* Set path to the directory
```
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

* Import library for aspell
```
go get github.com/trustmaster/go-aspell
```

* go run Program_name.go
```
Username : IP:Port
Friend IP : FriendIP+FriendPort
```

## Libraries used

* Aspell : used for spell check
* net    : used for tcp sockets
* bufio  : used to provide buffering and some help for textual I/O.
* time   : used for time strap

## Objective

* Make peer to peer connection 
* Make end-to-end encryption

## About the code

* Peer to peer connection is made using tcp sockets
* RSA algorithm is used for encryption and decryption
* XOR cipher is used for key exchange

## Concurrency

* Each peers runs a go routine for accepting a connection and anothor go routine to send the request
* For messageing two separate threads are created onr for read and one for write
* For Multiple users multiple threads are created, each having a read write pair





