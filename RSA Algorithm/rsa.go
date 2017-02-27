package main

import (
	//"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
        "bufio"
        //"strings"
)

func main() {

	// Generate RSA Keys
	myPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	myPublicKey := &myPrivateKey.PublicKey

	frPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}

	frPublicKey := &frPrivateKey.PublicKey

	fmt.Println("Private Key : ", myPrivateKey)
	fmt.Println("Public key ", myPublicKey)
	fmt.Println("Private Key : ", frPrivateKey)
	fmt.Println("Public key ", frPublicKey)

	//Encrypt  Message
       var send int
    var text1 string
    text1=""
    send=1
    fmt.Print("enter the message: ")
    for send==1 {
      reader := bufio.NewReader(os.Stdin)
      //fmt.Print("enter the message: ")
      text, _ := reader.ReadString('\n')
      text1=text1+text
      fmt.Print("Do you want to write more?")
      fmt.Scanf("%d",&send)
    }
	message := []byte(text1)
	label := []byte("")
	hash := sha256.New()
        
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, frPublicKey, message, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(" encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	

	// Decrypt Message
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, frPrivateKey, ciphertext, label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf(" decrypted [%x] to \n[%s]\n", ciphertext, plainText)

	

}
