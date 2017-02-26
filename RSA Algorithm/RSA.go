package main

 import (
   "crypto/rand"
   "fmt"
   "math/big"
 )

var e [256]big.Int
var d [256]big.Int
var temp [256]big.Int
var en [256]big.Int
var m [256]big.Int

func randomnumber() *big.Int {
   var p *big.Int
   var bits int
   var err error

   bits = 99

   p, err = rand.Prime(rand.Reader, bits)

   if err != nil {
      fmt.Println(err)
   }

   fmt.Printf("%d\n", p)
   return p
}


func ce(t,p,q *big.Int) {
   //var k,i *big.Int;
   //k=;
   k := big.NewInt(0)
   //var i *big.Int
   for i:=big.NewInt(2);i<t;i++ {
      if t%i==0 {
         continue;
      }

      flag := i.ProbablyPrime(4)
      //flag=prime(i);
      if (flag==true)&&(i!=p)&&(i!=q) {
        e[k]=i;
        var f big.Int
        f=cd(e[k]);
        if f>0 {
          d[k]=f;
          k++;
        }
        if k==255 {
         break;
        }
      }
   }
}

func cd(x *big.Int) *big.Int {
   var k big.Int
   k=1
   for {
      k=k+t;
      if k%x==0 {
         return k/x;
      }
   }
}

func encrypt() {
   var i,length,pt,k,key,ct *big.Int
   i=0
   key=e[0]
   length=len(m1)
   for (i!=length) {

      pt=m1[i]
      pt=pt-96
      k=1;
      for j:=0;j<key;j++ {
         k=k*pt
         k=k%n
      }
      temp[i]=k
      ct=k+96
      en[i]=ct
      i++
   }
   en[i]=-1
   fmt.Printf("\nTHE ENCRYPTED MESSAGE IS\n")
   for i:=0;en[i]!=-1;i++ {
      fmt.Printf("%c",en[i])
   }
}


func decrypt() {
   var pt,ct,key,k *big.Int
   i=0
   key=d[0]
   for en[i]!=-1 {
      ct=temp[i]
      k=1
      for j:=0;j<key;j++ {
        k=k*ct
        k=k%n
      }
      pt=k+96
      m[i]=pt
      i++

   }
   m[i]=-1
   fmt.Printf("\nTHE DECRYPTED MESSAGE IS\n")
   for i:=0;m[i]!=-1;i++ {
      fmt.Printf("%c",m[i])
   }
   fmt.Printf("\n")
   
}

 func main() {

   var p,q,n,t big.Int 
   //var arr []big.Int
     

   p=randomnumber()
   q=randomnumber()

   // if(p==q){
   //    q=randomnumber()

   // }
   n=pq
   t=(p-1)(q-1)

   var m1 string
   fmt.Printf("Enter message: ")
   fmt.Scanf("%s",&m1)
   for i:=0;i<len(m1);i++ {
      m[i]=m1[i]
   }
   ce(t,p,q)
   encrypt()
   decrypt()

 }