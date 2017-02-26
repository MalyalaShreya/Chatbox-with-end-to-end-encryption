package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"
import "time"
// import "strings"
// import "html/template"
// import "net/http"
var globalMap = make(map[string]bool)

func f(self string,self_port string) {

	var wg sync.WaitGroup
	wg.Add(1)
	go server(self,self_port)
	go AddFriend()

	wg.Wait();
}

func AddFriend() {

  client()
}

func server(self string,self_port string) {

    //s := make([]string, 0)
  //m := make(map[string]bool)
  	ln, _ := net.Listen("tcp", ":"+self_port)
     
    var i bool
    var user int
    user=0
    i=false
  	conn, _ := ln.Accept()
    i=true
    if i==true {
    var j int
    j=0
    fmt.Print("You got a friend request!\nDo you want to accept it?(Yes=1, No=0 )\n")
    fmt.Scanf("%d",&j) 
    if j==1 {
        user=user+1
        globalMap[conn.RemoteAddr().String()]=true
        
        // run loop forever (or until ctrl-c)
        var b bool
        //b=true
        //fmt.Println("Type :q! to Quit")
        for {
          // will listen for message to process ending in newline (\n)
          scanner := bufio.NewScanner(bufio.NewReader(conn))
          i:=0
          for scanner.Scan(){
             if scanner.Text()==""{
               break;
              }
          if i==0{
            fmt.Println("message received from client")
              i+=1
             }
            fmt.Println(scanner.Text())
            }
  
          reader := bufio.NewReader(os.Stdin)
          text:=""
          send:=1
          for send==1{
         fmt.Print("Text to send to client: ")
          newmessage, _ := reader.ReadString('\n')
          text=text+newmessage
         fmt.Print("Do you want to send?//if yes enter 1//if no enter 0??")
         fmt.Scanf("%d",&send)
    }
          // send new string back to client
          conn.Write([]byte(newmessage + "\n"))
        } 
      } else {
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

  fmt.Println("Type ")
	conn, _ := net.Dial("tcp", friendIP+":"+friendport)
  for {
      
      var send int
      var text1 string
    text1=""
    send=1
    for send==1 {
      reader := bufio.NewReader(os.Stdin)
      fmt.Print("Text to send to server: ")
      text, _ := reader.ReadString('\n')
      text1=text1+text
      // send to socket
      fmt.Print("Do you want to send?")
      fmt.Scanf("%d",&send)
    }
      fmt.Fprintf(conn, text + "\n")
      //Timer thread
      //Ask to cancel request or wait

      // listen for reply
      scanner := bufio.NewScanner(bufio.NewReader(conn))
     i:=0
    for scanner.Scan(){
     if scanner.Text()==""{
     break;
        }
    if i==0{
     fmt.Println("message received from server")
     i+=1
       }
    fmt.Println(scanner.Text())
   }
      timer := time.NewTimer(time.Second * 1)

    if message=="" {
      fmt.Println("Request not accepted\n")
      conn.Close()
      break
    } else {
        fmt.Print("Message from server: "+message)

    }

  }
}

func main() {

	var self,self_port string
	fmt.Print("Welcome to your chatbox!\nEnter your name: ")
	fmt.Scanf("%s",&self)
	fmt.Print("Enter your port: ")
	fmt.Scanf("%s",&self_port)
  var boolean bool
	boolean=true
  f(self,self_port)

  for boolean==true {

    fmt.Print("Options:\n1. Add friend\n2. Quit")
    var i int 
    fmt.Scanf("%d",&i)
    switch i {
   
      case 1:
        AddFriend()
      case 2:
        break

    }    
  }
	  
	fmt.Print("Bye!\nSee you soon :)");
  
}




