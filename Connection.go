package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"
import "github.com/trustmaster/go-aspell"
import "strings"
import "time"

type peer struct {

  IP string
  connection net.Conn
  //publicKey *rsa.PublicKey

}


type self struct {

  IP string

}


var Self self

var AllFriends map[string]peer
var PendingRequests map[string]peer


func Listen() {

  ln, _ := net.Listen("tcp", ":3000")

  for {

    conn, err := ln.Accept()
    if err != nil {
      fmt.Println(err.Error())
    }

    ip:=conn.RemoteAddr().String()
    PendingRequests[ip]=peer{ip,conn}
    //AllFriends[ip]=peer{ip,conn}

  }
}

func AcceptRequest() {

  if len(PendingRequests)!=0 {

    fmt.Print("Enter the IP of the person to accept the request: ")
    var ip string
    fmt.Scanf("%s",&ip)
    AllFriends[ip]=PendingRequests[ip]
    delete(PendingRequests,ip)
    fmt.Println(ip," is now your friend!\n")

  } else {
    fmt.Println("No pending requests!\n")
  }
  
}

func Delete() {

  if len(AllFriends)==0 {
    fmt.Println("You got no friends!\n")
  } else {

    fmt.Print("Enter the IP: ")
    var ip string
    fmt.Scanf("%s",&ip)
    if val,ok:=AllFriends[ip];ok {
      delete(AllFriends,ip)
      fmt.Println(val.IP,"is no longer your friend now.\n")

    } else {
      fmt.Println(ip, " was not your friend\n")
    }

  }

}

func chat(ip string) {
  fmt.Println("-----------------------------------------\n")  

  fmt.Print("To quit the chat with a friend type ':quit'\n")
  fmt.Println(".........................................\n")

  if val, ok := AllFriends[ip]; ok { 

    ch1:=make(chan int)
    ch2:=make(chan int)

    go func() {
      ch1<-Write(val.IP)
    }()
    go func (){
      ch2<-Read(val.IP)
    }()
    
    
    <-ch1
    <-ch2
    close(ch1)
    close(ch2)
    fmt.Println("---------------------------------------\n")  

  } else {

    fmt.Println(ip, " is not your friend.")

  }

}


func Read(ip string) (int) {

  for {

    
    buff:=make([]byte,0)
    tmp:=make([]byte,2048)
    n,_:=AllFriends[ip].connection.Read(tmp)
    buff=append(buff,tmp[:n]...)
    //AllFriends[ip].connection.Write([]byte("Seen"))
    current := time.Now()
    
    if(string(buff)==":quit\n") {
      fmt.Print(ip," went offline now\n")
      AllFriends[ip].connection.Write([]byte(":quit\n"))
      return 0
    }
    fmt.Print(ip," (",current.Format("2006-01-02 15:04:05"),"): ", string(buff))
  }

}

func Write(ip string) (int) {

  for {

    reader := bufio.NewReader(os.Stdin)
    newmessage, _ := reader.ReadString('\n')
    if newmessage==":quit\n" {
      AllFriends[ip].connection.Write([]byte(newmessage))
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
        } else {
            fmt.Print("[",v)
            fmt.Printf("] (Do you mean?): %s\n", strings.Join(speller.Suggest(v), ", "))
        }

    }

    AllFriends[ip].connection.Write([]byte(newmessage))
  }

}



func AddFriend() {

  var friendIP string
  var mutex = &sync.Mutex{}
  mutex.Lock()
  fmt.Print("Enter the IP: ")
  fmt.Scanf("%s",&friendIP)
  mutex.Unlock()

  conn, _ := net.Dial("tcp", friendIP)

  AllFriends[friendIP]=peer{friendIP,conn}

}

func showallfriends() {

  if len(AllFriends)==0 {
    fmt.Println("Forever alone! :(\n")
  } else {
    fmt.Println("Friends: ")
    fmt.Println("---------------------------------------\n")  

    for _,v:=range AllFriends {

      fmt.Println(v.IP)

    }
    fmt.Println()
  }

}

func showpendingrequests() {

  if len(PendingRequests)==0 {

    fmt.Println("No requests!\n")
  } else {

    fmt.Println("Pending requests: ")
    fmt.Println("---------------------------------------\n")  

    for _,v:=range PendingRequests {
      fmt.Println(v.IP)

    }
    fmt.Println()
  }

  
}

func main() {

  AllFriends = make(map[string]peer)
  PendingRequests=make(map[string]peer)
  
  var selfIP string
  Self = self{IP:""}
  
  go Listen()
  var i int

  var res bool
  res=true
    
    fmt.Println("\n")
    fmt.Print("                                                  *********************                                        \n")
    fmt.Print("                                                  *     MESSENGER     *                             \n")
    fmt.Print("                                                  *********************                               \n\n")
    fmt.Print("                                              Username :  ")
    fmt.Scanf("%s",&selfIP)
    Self.IP=selfIP
    fmt.Println("\n")

    fmt.Print("...................................................................................................................\n")
    fmt.Print(":  MyProfile                                                                      Home  Find Friends  *Active*    :\n")
    fmt.Print("...................................................................................................................\n\n")

    fmt.Print("--------------  --------------  --------------------  --------------------  ---------  -------------  -------------\n")
    fmt.Print("|1. AddFriend|  |2. MyFriends|  |3. Friend Requests|  |4. Confirm Request|  |5. Chat|  |6. Unfriend|  |7. Sign out|\n")
    fmt.Print("--------------  --------------  --------------------  --------------------  ---------  -------------  -------------\n\n")


    for res==true {
    fmt.Print("What's next?  ")

    fmt.Scanf("%d",&i)
    fmt.Println()

    switch (i) {

    case 1:
      AddFriend()
      break
    case 2:
      showallfriends()
      break
    case 3:
      showpendingrequests()
      break

    case 4:
      AcceptRequest()
      break
    case 5: 
      var ip string
      fmt.Println("Enter the IP of your friend: ")
      fmt.Scanf("%s",&ip)
      chat(ip)
      break
    case 6: 
      Delete()
      break

    case 7:
      res=false
      break

    }

  }

  fmt.Print("\nBye!\nSee you soon :)\n");
  
}




