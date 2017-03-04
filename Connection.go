package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"
import "crypto/x509" 
// import "bytes"

type peer struct {

  name string
  IP string
  port string
  connection net.Conn
  publicKey string

}

type self struct {

  name string
  IP string
  port string
  publicKey *rsa.PublicKey

}

var Self self
var SelfPrivateKey *rsa.PrivateKey

var AllFriends map[string]peer
var MyFriends map[string]peer

func Options() {

  var wg sync.WaitGroup
  wg.Add(1)
  go Listen()
  go AddFriend()
  // go GroupChat()
  wg.Wait()

}


func Listen() {

  ln, _ := net.Listen("tcp", ":"+Self.port)

  for {

    conn, err := ln.Accept()
    if err != nil {
      fmt.Println(err.Error())
    }

    var response string
    fmt.Println("You got a friend request!\nDo you want to accept it? (yes/no)")
    fmt.Scanf("%s",&response)
    if (response=="yes") {
      handleConn(conn)

    } else {
      fmt.Println("Request Declined")
    }

  }
}

func handleConn(conn net.Conn) {

  conn.Write([]byte(Self.name + "\n"))
  name, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("Message Received:", string(name))

  conn.Write([]byte(Self.IP + "\n"))
  ip, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("Message Received:", string(ip))

  conn.Write([]byte(Self.port + "\n"))
  port, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("Message Received:", string(port))

  Pub := Self.publicKey
  bytes, err := x509.MarshalPKIXPublicKey(Pub)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  str := string(bytes)

  conn.Write([]byte(str + "\n"))
  public, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("Message Received:", string(public))
  
  Friend := peer{name,ip,port,conn,public}
  AllFriends[name]=Friend
  MyFriends[name]=Friend

  chat(name)

}

func chat(name string) {

  var wg sync.WaitGroup
  wg.Add(2)
  go Read(name)
  go Write(name)
  wg.Wait()


}


func Read(name string) {

  for {
    message, _ := bufio.NewReader(AllFriends[name].connection).ReadString('\n')
    msg:=decrypt([]byte(message))
    fmt.Print("Message Received: ", msg)
  }
  

}

func Write(name string) {

  for {

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    newmessage, _ := reader.ReadString('\n')
    message:=encrypt(name,newmessage)
    AllFriends[name].connection.Write([]byte(message + "\n"))
  }
  

}

func encrypt(name string,msg string) (string){

  message := []byte(msg)
  label := []byte("")
  hash := sha256.New()
  key:=[]byte(AllFriends[name].publicKey)
  frPublicKey,err:= x509.ParsePKIXPublicKey(key)
  Public:=frPublicKey.(*rsa.PublicKey)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }


  ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, Public, message, label)

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

  SelfPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

  if err != nil {
    fmt.Println(err.Error)
    os.Exit(1)
  }

  myPublicKey := &SelfPrivateKey.PublicKey

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




