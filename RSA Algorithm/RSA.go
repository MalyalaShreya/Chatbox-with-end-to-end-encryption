//RSA in go
//Max prime  no. is 10000
//a bit slow
//We have to give an extra character at the end for termination
package main

import (
  "fmt"
  "math"
  "math/rand"
  "time"
)

const N = 10000
var Maxprimenumbers int
var e []int
var d []int
var temp []int
var en []int
var m []int
var is_prime [N]bool
var primes []int

// Only primes less than or equal to N will be generated

func random(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min) + min
}

func primenumbergeneration() {

  var x, y, n int
    nsqrt := math.Sqrt(N)

    

    for x = 1; float64(x) <= nsqrt; x++ {
        for y = 1; float64(y) <= nsqrt; y++ {
            n = 4*(x*x) + y*y
            if n <= N && (n%12 == 1 || n%12 == 5) {
                is_prime[n] = !is_prime[n]
            }
            n = 3*(x*x) + y*y
            if n <= N && n%12 == 7 {
                is_prime[n] = !is_prime[n]
            }
            n = 3*(x*x) - y*y
            if x > y && n <= N && n%12 == 11 {
                is_prime[n] = !is_prime[n]
            }
        }
    }

    for n = 5; float64(n) <= nsqrt; n++ {
        if is_prime[n] {
            for y = n * n; y < N; y += n * n {
                is_prime[y] = false
            }
        }
    }

    is_prime[2] = true
    is_prime[3] = true

    for x = 0; x < len(is_prime)-1; x++ {
        if is_prime[x] {
            primes = append(primes, x)
        }
    }

  Maxprimenumbers=len(primes)
    

}

func randomprimenumber() int {

  myrand := random(0,Maxprimenumbers)
  return primes[myrand]

}


func ce(t int,p int,q int) {
   var k,i int;
   k = 0
   for i=2;i<t;i++ {
      if t%i==0 {
         continue;
      }
      var flag bool
      flag = is_prime[i]
      if (flag==true)&&(i!=p)&&(i!=q) {
        e[k]=i;
        //e=append(e,i)
        var f int
        f=cd(e[k],t)
        if f>0 {
          flag=true
        } else {
          flag=false
        }
        if f>0 {
          d[k]=f;
          fmt.Print("d[",i,"]: ",d[k])
          k++;
        }
        if k==len(m) {
          fmt.Print("length of message: ",k)
         break;
        }
      }
   }
}

func cd(x int,t int) int {
   var k int
   k=1
   for {
      k=k+t;
      if k%x==0 {
         return k/x;
      }
   }
}

func encrypt(n int) {
   var i,length,pt,k,key,ct int
   i=0
   key=e[0]
   length=len(m)
   for (i!=length) {

      pt=m[i]
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
   fmt.Printf("e[%d]: -1",i)
   fmt.Printf("\nTHE ENCRYPTED MESSAGE IS\n")
   for i:=0;en[i]!=-1;i++ {
      fmt.Printf("%c",en[i])
   }
   fmt.Printf("Encrq")
}


func decrypt(n int) {
  fmt.Printf("Hi\n")
   var pt,ct,key,k,i int
   i=0
   key=d[0]
   fmt.Println("d[0]: ",d[0])
   for i<len(m)-1 {
    fmt.Printf("inside for\n")

      fmt.Printf("i= %d\n",i)
      ct=temp[i]
      k=1
      fmt.Print(key)
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

  var p,q,n,t int
  primenumbergeneration() 
     
  p=randomprimenumber()
  q=randomprimenumber()
  fmt.Println("p: ",p)
  fmt.Println("q: ",q)
  for p==q {
    q=randomprimenumber()
  }

  n=p*q
  t=(p-1)*(q-1)
  fmt.Println("t: ",t)
  fmt.Println("n: ",n)


   var m1 string
   fmt.Printf("Enter message: ")
   fmt.Scanf("%s",&m1)
   for i:=0;i<len(m1);i++ {
      m=append(m,int(m1[i]))
      e=append(e,0)
      d=append(d,0)
      en=append(en,0)
      temp=append(temp,0)
   }
  en=append(en,0)
  temp=append(temp,0)
   ce(t,p,q)
   encrypt(n)
   decrypt(n)

 }