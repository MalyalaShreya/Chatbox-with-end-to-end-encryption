package main

import (
    "github.com/andlabs/ui"
    "fmt"
)

func main() {
    err := ui.Main(func() {
        name := ui.NewEntry()
        ip:=ui.NewEntry()
        port:=ui.NewEntry()
        button := ui.NewButton("Add Friend")
        box := ui.NewVerticalBox()
        box.Append(ui.NewLabel("Enter name:"), false)
        box.Append(name, false)
        box.Append(ui.NewLabel("Enter IP address:"), false)
        box.Append(ip,false)
        box.Append(ui.NewLabel("Enter port"), false)
        box.Append(port,false)
        box.Append(button, false)
        window := ui.NewWindow("Add friend", 500, 500, false)
        s:=""
        

        name1 := ui.NewEntry()
        button1 := ui.NewButton("OK")
        message:=ui.NewEntry()
        messagebutton:=ui.NewButton("send")
        write:=ui.NewVerticalBox()
        box1 := ui.NewVerticalBox()
        box1.Append(ui.NewLabel("Enter name:"), false)
        box1.Append(name1, false)
        box1.Append(button1,false)
        box1.Append(ui.NewLabel("Enter message"), false)
        box1.Append(message,false)
        box1.Append(messagebutton, false)
        box1.Append(write,false)
        window1 := ui.NewWindow("chat", 500, 500, false)
       s1:=""
    
        button2:=ui.NewButton("Show Friends")
        show:=ui.NewVerticalBox()
        box2:=ui.NewVerticalBox()
        box2.Append(button2,false)
        box2.Append(show,false)
        window2:=ui.NewWindow("Friends list", 500, 500, false)
        s2:=""
        go func(){  
        window.SetChild(box)
        button.OnClicked(func(*ui.Button) {
             s=s+"hello"+name.Text()
            //greeting.SetText(s)
            fmt.Println("s",s)
        })
        }()       

        go func(){  
        window1.SetChild(box1)
        button1.OnClicked(func(*ui.Button) {
            s1="hello" 
        })
        }()

        go func(){  
        window2.SetChild(box2)
        button2.OnClicked(func(*ui.Button) {
               s2="hello" 
        })
        }()
        
       /* //box.Append(greeting,true)
        window1 := ui.NewWindow("Hello", 500, 500, false)
        window := ui.NewWindow("Hello", 500, 500, false)
        se:=""
        s:=""
        go func(){  
         
        window.SetChild(box)
        button.OnClicked(func(*ui.Button) {
        })
        }()
        go func(){
        
        window1.SetChild(newbox)
        send.OnClicked(func(*ui.Button){
           se=se+message.Text()
            fmt.Println("se",se)})
      }()*/
        
        
        window.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
         window1.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
         window2.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
        window.Show()
        window1.Show()
        window2.Show()
    })
    if err != nil {
        panic(err)
    }
}


