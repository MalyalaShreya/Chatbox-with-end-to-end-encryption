package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	conn := whatToDo()
	if conn == nil {
		log.Fatalln("Hag Diya")
	}

	go write(conn)
	read(conn)
}


func whatToDo() net.Conn {
	if os.Args[1] == "listen" {

		ln, err := net.Listen("tcp", os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		return conn
	}

	if os.Args[1] == "connect" {
		conn, err := net.Dial("tcp", os.Args[2])
		if err != nil {
			log.Fatalln(err)
		}

		return conn
	}

	return nil
}

func read(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print("MESSAGE RECIEVED : ", string(message))
	}
}

func write(conn net.Conn) {
	for {
		readMessage := bufio.NewReader(os.Stdin)
		readText, _ := readMessage.ReadString('\n')
		conn.Write([]byte(readText + "\n"))
	}
}