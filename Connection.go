//Chatbox with end-to-end encryption

//Self Port is kept constant and equal to 3000

package main

import "net"
import "fmt"
import "bufio"
import "os"
import "sync"
import "github.com/trustmaster/go-aspell"
import "strings"
import "time"
import "math/rand"
import "math/big"



//Struct to store the details of friend
type peer struct {
	IP         string   //IP is unique key
	connection net.Conn // connection stores the value of net.Conn when connection is made
	//(n,e) tuple stores the public key of friend
	n int
	e int
}

// Struct to store self details
type self struct {
	IP string
}

//Struct to store the prime numbers used for encryption decryption
type prime struct {
	p int
	q int
}

//Struct to store self public key
type publickey struct {
	n int
	e int
}

//Struct to store self private key
type privatekey struct {
	n int
	d int
}


//Initializing self public and private keys
var PublicKey publickey
var PrivateKey privatekey

var Self self //Initilalizing self struct

var AllFriends map[string]peer      //Map to store list of all friends
var PendingRequests map[string]peer //Map to store list of pending friend requests
var quit_read, quit_write int



//Function to listen all the friend requets
func Listen() {

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
		return
	}
	KeyGenerate() //Generates public and private keys

	//Making the connection with peer
	for {

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		ip := conn.RemoteAddr().String()           //Stores the IP of the peer in the form of string
		x, y := ReadKey(conn)                      // Reads the key sent from peer
		WriteKey(conn)                             // Sends self public and private key to peer
		PendingRequests[ip] = peer{ip, conn, x, y} //Stores all the details of the peer in pending queue

	}
}



//Function to accept friend request
func AcceptRequest() {

	if len(PendingRequests) != 0 { // If there is at least one friend request

		fmt.Print("Enter the IP of the person to accept the request: ")
		var ip string
		fmt.Scanf("%s", &ip)
		AllFriends[ip] = PendingRequests[ip] //Add the friend to your friend list
		delete(PendingRequests, ip)
		fmt.Println(ip, " is now your friend!\n")

	} else {
		fmt.Println("No pending requests!\n")
	}

}



//Blocking a person
func Delete() {

	if len(AllFriends) == 0 {
		fmt.Println("You got no friends!\n")
	} else { //If there is at least one friend

		fmt.Print("Enter the IP: ")
		var ip string
		fmt.Scanf("%s", &ip)
		if val, ok := AllFriends[ip]; ok {
			delete(AllFriends, ip) // Deleting the person from friend list
			fmt.Println(val.IP, "is no longer your friend now.\n")

		} else {
			fmt.Println(ip, " was not your friend\n") //If he/she was not your friend
		}

	}

}



//Function to beging the chat with your peer
func chat(ip string) {
	fmt.Println("-----------------------------------------\n")

	fmt.Print("To quit the chat with a friend type ':quit'\n")
	fmt.Println(".........................................\n")

	//If the friend is in your friendlist
	if val, ok := AllFriends[ip]; ok {

		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			ch1 <- Write(val.IP)
		}()
		go func() {
			ch2 <- Read(val.IP)
		}()

		<-ch1
		<-ch2
		close(ch1)
		close(ch2)
		fmt.Println("---------------------------------------\n")
		//delete(AllFriends,ip)

	} else {

		fmt.Println(ip, " is not your friend.")

	}

}



//Function to read messages
func Read(ip string) int {

	for {

		buff := make([]byte, 0)
		tmp := make([]byte, 2048)
		n, _ := AllFriends[ip].connection.Read(tmp)
		buff = append(buff, tmp[:n]...)
		decrypted := decrypt(buff) //Recieves decrypted messages
		current := time.Now()      //To show the time when message is recieved

		//If the friend quits
		if decrypted == ":quit\n" {
			fmt.Print(ip, " went offline now\n")
			AllFriends[ip].connection.Write(encrypt(":quit\n", AllFriends[ip].n, AllFriends[ip].e))
			return 0
		}
		fmt.Print(ip, " (", current.Format("2006-01-02 15:04:05"), "): ", decrypted)
	}

}



//Function to send messages
func Write(ip string) int {

	for {

		if val, ok := AllFriends[ip]; ok { //Checks if the friend is in your friend list

			reader := bufio.NewReader(os.Stdin)
			newmessage, _ := reader.ReadString('\n')
			if newmessage == ":quit\n" { //For quitting the chat

				AllFriends[val.IP].connection.Write(encrypt(newmessage, AllFriends[ip].n, AllFriends[ip].e))
				return 0
			}
			//Words contains the array of words obtained from newmessage
			words := strings.Fields(newmessage)

			//For spell checker
			speller, err := aspell.NewSpeller(map[string]string{
				"lang": "en_US",
			})
			if err != nil {
				fmt.Errorf("Error: %s", err.Error())
				return 0
			}

			for _, v := range words {

				//Checks if the words are correct or not
				if speller.Check(v) {
				} else {
					fmt.Print("[", v)
					//If the words are not correct then suggests the most appropriate word
					fmt.Printf("] (Do you mean?): %s\n", strings.Join(speller.Suggest(v), ", "))
				}

			}

			//Prints the encrypted text
			AllFriends[ip].connection.Write(encrypt(newmessage, AllFriends[ip].n, AllFriends[ip].e))
		} else {
			return 0
		}
	}

}



//Function to add the friend
func AddFriend() {

	var friendIP string
	var mutex = &sync.Mutex{}
	mutex.Lock()
	fmt.Print("Friend IP: ")
	fmt.Scanf("%s", &friendIP)
	mutex.Unlock()

	conn, err := net.Dial("tcp", friendIP) //Dial is similar to sending the friend request
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
		return
	}
	//After the connection is made Writekey and Readkey are used to send and read keys
	WriteKey(conn)
	x, y := ReadKey(conn)
	AllFriends[friendIP] = peer{friendIP, conn, x, y} //Adds the person in friend list

}



// Function to display all the friends in the friend list
func showallfriends() {

	if len(AllFriends) == 0 {
		fmt.Println("Forever alone! :(\n")
	} else {
		fmt.Println("Friends: ")
		fmt.Println("---------------------------------------\n")

		for _, v := range AllFriends {

			fmt.Println(v.IP)

		}
		fmt.Println()
	}

}



//Function to display all the pending friend requests
func showpendingrequests() {

	if len(PendingRequests) == 0 {

		fmt.Println("No requests!\n")
	} else {

		fmt.Println("Pending requests: ")
		fmt.Println("---------------------------------------\n")

		for _, v := range PendingRequests {
			fmt.Println(v.IP)

		}
		fmt.Println()
	}

}



//Function to read the keys
func ReadKey(conn net.Conn) (int, int) {
	buff := make([]byte, 0)
	tmp := make([]byte, 2048)
	n, err := conn.Read(tmp)
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
		return 0, 0
	}
	buff = append(buff, tmp[:n]...)
	return int(buff[0] ^ 11), int(buff[1] ^ 11) //The key is first decrypted nad then stored

}



//Function to send the key
func WriteKey(conn net.Conn) {
	arr1 := make([]byte, 0)
	//The key is first encrypted and then sent

	//Keys are sencrypted and decrypted using XOR
	//We know, a^b^b=a so a^b is sent to peer and it recieves the key and XOR it wih b to get the key
	arr1 = append(arr1, byte(PublicKey.n^11))
	arr1 = append(arr1, byte(PublicKey.e^11))
	conn.Write(arr1)

	return
}



//Function to get the GCD of two numbers
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}



//Function to generate keys for encryption and decryption
func KeyGenerate() {
	//p*q should be in the range of 200 to 256 for better encrytion
	var arr [12]prime
	arr[0].p = 11
	arr[0].q = 13
	arr[1].p = 11
	arr[1].q = 17
	arr[2].p = 11
	arr[2].q = 19
	arr[3].p = 11
	arr[3].q = 23
	arr[4].p = 13
	arr[4].q = 17
	arr[5].p = 13
	arr[5].q = 19

	arr[6].p = 13
	arr[6].q = 11
	arr[7].p = 17
	arr[7].q = 11
	arr[8].p = 19
	arr[8].q = 11
	arr[9].p = 23
	arr[9].q = 11
	arr[10].p = 17
	arr[10].q = 13
	arr[11].p = 19
	arr[11].q = 13

	rand.Seed(time.Now().Unix())           //Generating random numbers
	x := rand.Int() % 12                   // Genrating random number from 0 to 11
	n := arr[x].p * arr[x].q               //n=p*q
	phi := (arr[x].p - 1) * (arr[x].q - 1) //phi=(p-1)*(q-1)
	epos := make([]int, 0)
	for i := 1; i < phi; i++ {
		if gcd(i, phi) == 1 {
			epos = append(epos, i) //all possible values of e such that e is the coprime of phi
		}
	}
	m := rand.Int() % len(epos) //selecting e randomly form all the possible values of e
	e := epos[m]
	d := 0
	s := 0
	//generating the value d such that ed=1 mod phi
	for s != 1 {
		s = (d * e) % phi
		d++
	}
	d--

	//Storing self public and private keys
	PublicKey.n = n
	PublicKey.e = e
	PrivateKey.n = n
	PrivateKey.d = d
}



//For encyption and decryption RSA algorithm is followed
//Function to encrypt the text

func encrypt(msg string, n int, e int) []byte {
	msg1 := []byte(msg)
	cipher := make([]byte, 0)

	//(n,e) tuple is public keys of friend
	//for every byte of message we encrypt the value of the byte using the formula c=m^e mod n
	for i := range msg1 {
		c1 := new(big.Int)
		c1 = c1.Exp(big.NewInt(int64(msg1[i])), big.NewInt(int64(e)), nil)
		c2 := new(big.Int)
		c2 = c2.Mod(c1, big.NewInt(int64(n)))
		c3 := c2.Int64()
		//append all the encrypted bytes of message into a slice which gives the complete encrypted message of the plain text
		cipher = append(cipher, byte(c3))
	}
	return cipher
}



//Function to decrypt the text
//(n,d) tuple is the private key
//decrypt every byte of the encrypted message using the private key using the formula m=c^d mod n

func decrypt(cipher []byte) string {
	msg1 := make([]byte, 0)

	for i := range cipher {
		c1 := new(big.Int)
		c1 = c1.Exp(big.NewInt(int64(cipher[i])), big.NewInt(int64(PrivateKey.d)), nil)
		c2 := new(big.Int)
		c2 = c2.Mod(c1, big.NewInt(int64(PrivateKey.n)))
		c3 := c2.Int64()
		//append all the decrypted bytes of ciphertext into a slice which gives the complete plain text
		msg1 = append(msg1, byte(c3))
	}

	msg := string(msg1)
	return msg
}



//Main function from where the program execution begins
func main() {

	//Initializing maps
	AllFriends = make(map[string]peer)
	PendingRequests = make(map[string]peer)

	var selfIP string
	Self = self{IP: ""}

	go Listen() //go routine to listen
	var i int

	var res bool
	res = true

	fmt.Println("\n")
	fmt.Print("                                                  *********************                                        \n")
	fmt.Print("                                                  *     MESSENGER     *                             \n")
	fmt.Print("                                                  *********************                               \n\n")
	fmt.Print("                                              Username :  ")
	fmt.Scanf("%s", &selfIP)
	Self.IP = selfIP
	fmt.Println("\n")

	fmt.Print("...................................................................................................................\n")
	fmt.Print(":  MyProfile                                                                      Home  Find Friends  *Active*    :\n")
	fmt.Print("...................................................................................................................\n\n")

	fmt.Print("--------------  --------------  --------------------  --------------------  ---------  ----------  -------------\n")
	fmt.Print("|1. AddFriend|  |2. MyFriends|  |3. Friend Requests|  |4. Confirm Request|  |5. Chat|  |6. Block|  |7. Sign out|\n")
	fmt.Print("--------------  --------------  --------------------  --------------------  ---------  ----------  -------------\n\n")

	for res == true {
		fmt.Print("What's next?  ")

		fmt.Scanf("%d", &i)
		fmt.Println()

		switch i {

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
			fmt.Print("Enter the IP of your friend: ")
			fmt.Scanf("%s", &ip)
			chat(ip)
			break
		case 6:
			Delete()
			break

		case 7:
			res = false
			break

		}

	}
	fmt.Print("                                                  *********************                                        \n")
	fmt.Print("                                                  *      SAYONARA!    *                                        \n")
	fmt.Print("                                                  *********************                                        \n")

}
