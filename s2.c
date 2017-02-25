#include<stdio.h>
#include<stdlib.h>
#include<math.h>
#include<time.h>


// declairing function prototypes....
long long int power(long long int x, unsigned long long int y, long long int p);
int miillerTest(long long int d, long long int n);
int isPrime(long long int n, int k);

// start of main() function....
int main()
{
	int k = 32;  // Number of iterations
 
	unsigned long long int r,min,max,n,s;
	time_t t;
	srand((unsigned) time(&t));

	min = pow(2,k-1);
	max = pow(2,k);

	//printf("hi\n");
	r= 100 * ( (int)(log(k)/log(2)) +1);
	long long int d=1;	// a variable to check condition for while loop
	s=r;		// another variable to reset the value of r.
	
	// start of outer while lopp....
	// this loop is to iterate again and again until and unless the prime number is found.
	while(d!=n)
	{
		d=0;
		r=s;
		//printf("yo\n");
		// start of inner while loop....
		while(r>0)
		{
			n = rand() % max + min;		// generating the random number.
			r = r-1;
			//printf("%lld : %lld\n",r,n);			
			if(isPrime(n,k) == 1)
			{
				d=n;
				break;
			}
		} // end of inner while loop.
	} // end of outer while loop.
	//printf("yocum\n");
	
	printf("%lld",d);
 
        return 0;
}
// end of main() function....


// start of power() function.... 
// Utility function to do modular exponentiation. It returns (x^y) % p
long long int power(long long int x, unsigned long long int y, long long int p)
{
    long long int res = 1;      // Initialize result
    x = x % p;  // Update x if it is more than or equal to p
    while (y > 0)
    {
        // If y is odd, multiply x with result
        if (y & 1)
            res = (res*x) % p;
 
        // y must be even now
        y = y>>1; // y = y/2
        x = (x*x) % p;
    }
    return res;
}
// end of power() function.... 

// start of millerTest() function.....
// This function is called for all k trials. It returns false if n is composite and returns false if n is probably prime.
// d is an odd number such that  d*2<sup>r</sup> = n-1 for some r >= 1
int miillerTest(long long int d, long long int n)
{
    // Pick a random number in [2..n-2]
    // Corner cases make sure that n > 4
    long long int a = 2 + rand() % (n - 4);
 
    // Compute a^d % n
    long long int x = power(a, d, n);
 
    if (x == 1  || x == n-1)
       return 1;
 
    // Keep squaring x while one of the following doesn't happen
    // (i)   d does not reach n-1
    // (ii)  (x^2) % n is not 1
    // (iii) (x^2) % n is not n-1
    while (d != n-1)
    {
        x = (x * x) % n;
        d *= 2;
 
        if (x == 1)
	      return -1;
        if (x == n-1)    
		return 1;
    }
 
    // Return composite
    return -1;
}
// end of millerTest() function....

// start of isPrime() function.... 
// It returns false if n is composite and returns true if n is probably prime.  k is an input parameter that determines accuracy level.
// Higher value of k indicates more accuracy.
int isPrime(long long int n, int k)
{
    // Corner cases
    if (n <= 1 || n == 4)  
		return -1;
    if (n <= 3) 
		return 1;
 
    // Find r such that n = 2^d * r + 1 for some r >= 1
    long long int d = n - 1;
    while (d % 2 == 0)
        d /= 2;
 
    // Iterate given nber of 'k' times
    for (int i = 0; i < k; i++)
         if (miillerTest(d, n) == -1)
              return -1;
 
    return 1;
}
// end of isPrime() function.... 

