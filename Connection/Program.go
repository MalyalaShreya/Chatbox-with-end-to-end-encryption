package main

//worked with Matt Pozderac

import (
	"fmt"
	//"os"
	//"controlpkg"
    //"strconv"
    "html/template"
    "log"
    "net/http"
    //"strings"
    // "errors"
    // "net"
    // "os"
)


func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("Output.html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        // logic part of log in
        fmt.Println(r.Form)
        fmt.Println("Myself:", r.Form["myself"])
        //fmt.Println("ip1:", strconv.Atoi(r.FormValue("ip1")))
        fmt.Println("Friend:", r.Form["username"])
        fmt.Println("ip2:", r.Form["ip2"])
        var s string
        s=r.Form["ip2"]
        l:=s[1:len(s)-1]
        fmt.Println("IPString:",l)
        fmt.Println("MyText:", r.Form["mytext"])
        fmt.Println("FriendText:", r.Form["friendtext"])
        fmt.Fprintf(w, "MyText: %s\nFriendText: %s\n", r.Form["mytext"],r.Form["friendtext"])

    }
}

// func GetMyIP() (IP string) {
//     ip, err := myExternalIP()
//     if err != nil {
//         log.Println("could not get my external adress!")
//         HandleErr(err)
//     } else {
//         log.Printf("myExternalAdress = %s", ip)
//     }
//     IP = ip //external IP
//     return
// }

func main() {

    //http.HandleFunc("/", sayhelloName) // setting router rule
    http.HandleFunc("/login", login)
    err := http.ListenAndServe(":3000", nil) // setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

}


