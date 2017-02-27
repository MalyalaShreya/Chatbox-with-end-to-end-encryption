package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"

var allClients map[*Client]int
var allServers map[*Server]int

type Client struct {
    incoming chan string
    ipandport    string
    outgoing   chan string
    // reader     *bufio.Reader
    // writer     *bufio.Writer
    conn       net.Conn
    connection *Client
}

type Server struct {
    incoming chan string
    outgoing chan string
    ipandport  string
    name       string
    // reader     *bufio.Reader
    // writer     *bufio.Writer
    conn       net.Conn
    connection *Server

}



func (client *Client) Read() {
    for {
        reader := bufio.NewReader(os.Stdin)
        line, err := reader.ReadString('\n')
        fmt.Println("Text to client: ",line)
        if line==":q" {
          client.conn.Close()
          fmt.Println("Closed")
          break
        } else if line==":delete" {
          client.conn.Close()
          delete(allClients, client)
          if client.connection != nil {
            client.connection = nil

          }
          client = nil
          fmt.Println("Deleted")


        } else if err == nil {
            // if client.connection != nil {
            //     client.connection.outgoing <- line
            // send new string back to client
            client.conn.Write([]byte(line + "\n"))
            
        } else {
            break
        }

    }

    client.conn.Close()
    delete(allClients, client)
    if client.connection != nil {
        client.connection = nil
    }
    client = nil
}

func (server *Server) Read() {
    for {
        
        reader := bufio.NewReader(os.Stdin)
        line, err := reader.ReadString('\n')
        fmt.Println("Text to server: ",line)
        if line==":q" {
          server.conn.Close()
          fmt.Println("Closed")

          break
        } else if line==":delete" {
          server.conn.Close()
          delete(allServers, server)
          if server.connection != nil {
            server.connection = nil

          }
          server = nil
          fmt.Println("Deleted")


        } else if err == nil {
            // if server.connection != nil {
            //     server.connection.outgoing <- line
            // send new string back to client
            server.conn.Write([]byte(line + "\n"))
            
        } else {
            break
        }

    }

    server.conn.Close()
    delete(allServers, server)
    if server.connection != nil {
        server.connection = nil
    }
    server = nil
}

func (client *Client) Write() {
    // for data := range client.incoming {
    //     fmt.Println("Text from client: ",data)

    //     if data==":q" {
    //       client.conn.Close()
    //       fmt.Println("Closed")

    //       break
    //     }
    //     if data==":delete" {
    //       client.conn.Close()
    //       delete(allClients, client)
    //       if client.connection != nil {
    //         client.connection.connection = nil

    //       }
    //       client = nil
    //       fmt.Println("Deleted")


    //     }
    //     client.conn.Write([]byte(data + "\n"))
    //     //client.writer.Flush()
    // }
  for {
    message, _ := bufio.NewReader(client.conn).ReadString('\n')
            // output message received
    fmt.Print("Message Received from server:", string(message))
  }
}
func (server *Server) Write() {
    // for data := range server.incoming {
    //     fmt.Println("Text from server: ",data)

    //     if data==":q" {
    //       server.conn.Close()
    //       fmt.Println("Closed")
    //       break
    //     }
    //     if data==":delete" {
    //       server.conn.Close()
    //       delete(allServers, server)
    //       if server.connection != nil {
    //         server.connection.connection = nil
    //       }
    //       server = nil
    //       fmt.Println("Deleted")

    //     }
    //     server.conn.Write([]byte(data + "\n"))
    //     //server.writer.Flush()
    // }
  for {

    message1, _ := bufio.NewReader(server.conn).ReadString('\n')
            // output message received
    fmt.Print("Message Received from client:", string(message1))
  }

}

func (client *Client) Listen() {
    go client.Read()
    go client.Write()
}

func (server *Server) Listen() {
    go server.Read()
    go server.Write()
}
func f(self string,self_port string) {

	var wg sync.WaitGroup
	wg.Add(2)
	go server(self,self_port)
	go client()

	wg.Wait();
}

func AddFriend() {
  client()
}
func NewClient(connection net.Conn) *Client {
    // writer := bufio.NewWriter(connection)
    // reader := bufio.NewReader(connection)

    client := &Client{
        incoming: make(chan string),
        outgoing: make(chan string),
        conn:     connection,
        // reader:   reader,
        // writer:   writer,
    }
    client.Listen()

    return client
}

func NewServer(connection net.Conn) *Server {
    //fmt.Println("Newserver")
    // writer := bufio.NewWriter(connection)
    // reader := bufio.NewReader(connection)

    server := &Server{
        incoming: make(chan string),
        outgoing: make(chan string),
        conn:     connection,
        // reader:   reader,
        // writer:   writer,
    }
    server.Listen()

    return server
}



func server(self string,self_port string) {

  	listener, _ := net.Listen("tcp", ":"+self_port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err.Error())
        }
        client := NewClient(conn)
        fmt.Println("You got a friend request!\nDo you want to accept it?\n(yes/no)")
        var answer string
        fmt.Scanf("%s",&answer)
        if answer=="yes" {
          allClients[client] = 1
          client.ipandport=conn.RemoteAddr().String()
          fmt.Println("length: ",len(allClients))

          for clientList, _ := range allClients {
              if clientList.connection == nil {
                  client.connection = clientList
                  clientList.connection = client
                  fmt.Println("Connected")
              }

              fmt.Println("inside for")
          }
          fmt.Println("length: ",len(allClients))
        } else {
          fmt.Println("length: ",len(allClients))
          //delete(allClients,client)
          client.conn.Close()
          delete(allClients, client)
          if client.connection != nil {
              client.connection.connection = nil
          }
          client = nil
          fmt.Print("Sorry for bad experience :(\n")
        }
        
    }
    
    fmt.Print("Back to server\n")
  

}

func client() {

	var friend,friendIP,friendport string
  var mutex = &sync.Mutex{}
  mutex.Lock()
	fmt.Print("Enter the name of your friend: ")
	fmt.Scanf("%s",&friend)
	fmt.Print("Enter the IP: ")
	fmt.Scanf("%s",&friendIP)
	fmt.Print("Enter the port: ")
	fmt.Scanf("%s",&friendport)
  mutex.Unlock()


	conn, _ := net.Dial("tcp", friendIP+":"+friendport)
  server := NewServer(conn)
  server.name=friend
  server.ipandport=friendIP+":"+friendport
  allServers[server] = 1
  for serverList, _ := range allServers {
    if serverList.connection == nil {
      server.connection = serverList 
      serverList.connection = server
      }
  }
}

func main() {

  fmt.Print("      ---------------------------------------      \nBasic commands\n1. To quit the chat with a friend type :q\n2. To delete a friend type :delete\n      ---------------------------------------      \n ")
	var self,self_port string
  allClients = make(map[*Client]int)
  allServers = make(map[*Server]int)

	fmt.Print("Welcome to your chatbox!\nEnter your name: ")
	fmt.Scanf("%s",&self)
	fmt.Print("Enter your port: ")
	fmt.Scanf("%s",&self_port)
  //var boolean bool
	//boolean=true
  // client()

  f(self,self_port)

  // for boolean==true {

  //   fmt.Print("Options:\n1. Add friend\n2. Quit")
  //   var i int 
  //   fmt.Scanf("%d",&i)
  //   switch i {
   
  //     case 1:
  //       AddFriend()
  //     case 2:
  //       break

  //   }    
  // }
	  
	fmt.Print("Bye!\nSee you soon :)\n");
  
}




