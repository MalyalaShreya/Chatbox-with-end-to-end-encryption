package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"
import "github.com/trustmaster/go-aspell"
import "strings"
import "time"
import "encoding/gob"


type peer struct {

  name string
  IP string
  port string
  connection net.Conn
  publicKey *rsa.PublicKey

}


type recieve struct {

  name string
  IP string
  port string
  publicKey *rsa.PublicKey

}

type self struct {

  name string
  IP string
  port string
  publicKey *rsa.PublicKey

}

type request struct {

  IP string
  port string
  conn net.Conn

}

var Self self
var SelfPrivateKey *rsa.PrivateKey

var AllFriends map[string]peer
var MyFriends map[string]peer
var PendingRequests map[string]request


func Options() {

  var wg sync.WaitGroup
  wg.Add(1)
  //go Listen()
  go AddFriend()
  wg.Wait()

}


func Listen() {

  ln, _ := net.Listen("tcp", ":"+Self.port)

  for {

    conn, err := ln.Accept()
    if err != nil {
      fmt.Println(err.Error())
    }

    ip, port, err := net.SplitHostPort(conn.RemoteAddr().String())

    PendingRequests[ip]=request{ip,port,conn}

    var response string
    fmt.Println("You got a friend request!\nDo you want to accept it? (yes/no)")
    fmt.Scanf("%s",&response)
    if (response=="yes") {

      handleConn(conn)

    } else {

      
      //delete(sessions, "moo")
      delete(PendingRequests,ip)

      fmt.Println("Request Declined")
    }

  }
}

func handleConn(conn net.Conn) {

  ch:=make(chan string)

  go func() {
    encode(conn)
  }()
  go func (){
    ch<-decode(conn) 
  }()
  
  //var name string
  <-ch
  close(ch)
  fmt.Println("Done")

  

}

//chat(name)

func encode(conn net.Conn) {

  fmt.Println("Encoder")
  str := recieve{Self.name,Self.IP,Self.port,Self.publicKey}
  encoder := gob.NewEncoder(conn)
  encoder.Encode(str)
  fmt.Println("Encoder2")

}

func decode(conn net.Conn) (string) {
  fmt.Println("Decoder1")
  var person recieve
  decoder := gob.NewDecoder(conn)
  decoder.Decode(&person)
  AllFriends[person.name]=peer{person.name,person.IP,person.port,conn,person.publicKey}
  MyFriends[person.name]=AllFriends[person.name]
  fmt.Println("Decoder2")
  return person.name

}


func chat(name string) {

  ch1:=make(chan int)
  ch2:=make(chan int)

  go func() {
    ch1<-Write(name)
  }()
  go func (){
    ch2<-Read(name)
  }()
  
  
  <-ch1
  <-ch2
  close(ch1)
  close(ch2)
  fmt.Println("---------------------------------------\n")  


}


func Read(name string) (int) {

  for {
    message, _ := bufio.NewReader(AllFriends[name].connection).ReadString('\n')
    msg:=decrypt([]byte(message))
    fmt.Println("Message recieved at: ",time.Now().Format(time.RFC850))
    fmt.Print("Message Received: ", msg)
    if(msg==":quit") {
      return 0
    }
  }
  


}

func Write(name string) (int) {

  for {

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    newmessage, _ := reader.ReadString('\n')
    if newmessage==":quit" {
      message:=encrypt(name,newmessage)
      AllFriends[name].connection.Write([]byte(message + "\n"))
      return 0
    }

    words := strings.Fields(newmessage)

    speller, err := aspell.NewSpeller(map[string]string{
        "lang": "en_US",
    })
    if err != nil {
        fmt.Errorf("Error: %s", err.Error())
        return 0
    }

    for _,v := range words {

        if speller.Check(v) {
            //fmt.Print("OK\n")
        } else {
            fmt.Print(v)
            fmt.Printf("(suggestions): %s\n", strings.Join(speller.Suggest(v), ", "))
        }

    }

    message:=encrypt(name,newmessage)
    fmt.Println("Msg sent at: ",time.Now().Format(time.RFC850))
    AllFriends[name].connection.Write([]byte(message + "\n"))
  }

}

func encrypt(name string,msg string) (string){

  message := []byte(msg)
  label := []byte("")
  hash := sha256.New()
  ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, AllFriends[name].publicKey, message, label)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  fmt.Printf(" encrypted [%s] to \n[%x]\n", string(message), ciphertext)
  fmt.Println()

  return string(ciphertext)
}

func decrypt(msg []byte) (string){


  label := []byte("")
  hash := sha256.New()
  plainText, err := rsa.DecryptOAEP(hash, rand.Reader,SelfPrivateKey,msg,label)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  fmt.Printf(" decrypted [%x] to \n[%s]\n", msg, plainText)
  return string(plainText)

}

func AddFriend() {

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
  handleConn(conn)
}

func main() {

  AllFriends = make(map[string]peer)
  MyFriends = make(map[string]peer)
  PendingRequests = make(map[string]request)

  SelfPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

  if err != nil {
    fmt.Println(err.Error)
    os.Exit(1)
  }

  myPublicKey := &SelfPrivateKey.PublicKey
  fmt.Println("my private key",SelfPrivateKey)
  fmt.Println("my publickey",myPublicKey)
  fmt.Print("      ---------------------------------------      \nBasic commands\n1. To quit the chat with a friend type :q\n2. To delete a friend type :delete\n      ---------------------------------------      \n ")
  var name,self_port,selfIP string
  Self = self{name:"" ,IP:"",port:"",publicKey:myPublicKey}


  fmt.Print("Welcome to your chatbox!\nEnter your name: ")
  fmt.Scanf("%s",&name)
  fmt.Print("Enter your public IP: ")
  fmt.Scanf("%s",&selfIP)
  fmt.Print("Enter your port: ")
  fmt.Scanf("%s",&self_port)
  
  Self.name=name  
  Self.IP=selfIP
  Self.port=self_port

  Options()
  
  fmt.Print("Bye!\nSee you soon :)");
  
  
}




