#include <sys/types.h>
#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <netdb.h>
#include <string.h>
#include <pthread.h>
//#include <mutex>

pthread_mutex_t mutex,mutex1;
char self[50];
int self_port;
char IP[20];
char Name[50];
int port;

void error(char*msg)
{
	perror(msg);
 	exit(1);
}

void* server()
{
  printf("Hi %s, This is your chat room\n",self);
	pthread_mutex_t mutex;
	//pthread_mutex_lock(&mutex);
	int sock, connected, bytes_recieved , true = 1;  
        char send_data [1024] , recv_data[1024];       

        struct sockaddr_in server_addr,client_addr;    
        int sin_size;
        
        if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            perror("Socket");
            exit(1);
        }

        if (setsockopt(sock,SOL_SOCKET,SO_REUSEADDR,&true,sizeof(int)) == -1) {
            perror("Setsockopt");
            exit(1);
        }
        
        server_addr.sin_family = AF_INET;         
        server_addr.sin_port = htons(self_port);     
        server_addr.sin_addr.s_addr = INADDR_ANY; 
        bzero(&(server_addr.sin_zero),8); 

        if (bind(sock, (struct sockaddr *)&server_addr, sizeof(struct sockaddr))
                                                                       == -1) {
            perror("Unable to bind");
            exit(1);
        }

        if (listen(sock, 5) == -1) {
            perror("Listen");
            exit(1);
        }
		
        printf("Waiting for your friends to come online......\n");
        fflush(stdout);
        int boolean=1;
        while(boolean)
        { 
          sin_size = sizeof(struct sockaddr_in);

          if(accept(sock, (struct sockaddr *)&client_addr,(socklen_t *)&sin_size)>=0)
          {
            printf("\n I got a connection from (%s , %d)",
                         inet_ntoa(client_addr.sin_addr),ntohs(client_addr.sin_port));

                  while (1)
                  {
                    printf("\n SEND (q or Q to quit) : ");
                    scanf("%s",send_data);
                    
                    if (strcmp(send_data , "q") == 0 || strcmp(send_data , "Q") == 0)
                    {
                      send(connected, send_data,strlen(send_data), 0); 
                      close(connected);
                      break;
                    }
                     
                    else
                       send(connected, send_data,strlen(send_data), 0);  

                    bytes_recieved = recv(connected,recv_data,1024,0);

                    recv_data[bytes_recieved] = '\0';

                    if (strcmp(recv_data , "q") == 0 || strcmp(recv_data , "Q") == 0)
                    {
                      close(connected);
                      break;
                    }

                    else 
                    printf("\n RECIEVED DATA = %s " , recv_data);
                    fflush(stdout);
                  }

          }
          if(1)
          {
            printf("1. Add friend\n2. Leave chat room\n");
            int i;
            scanf("%d",&i);
            switch(i)
            {
              case 1:
              {
                
                printf("Enter the name of your friend: ");
                scanf("%s",Name);
                printf("Enter the IP address of your friend: ");
                scanf("%s",IP);
                printf("Enter the port: ");
                scanf("%d",&port);
                client();
                break;
              }

              default:
              boolean=0;
              //Do Nothing
            }
          }
        }
          
      close(sock);
}

int client()
{
  printf("Hello %s,you are now connected to %s\n",self,Name);
	int sock, bytes_recieved;  
        char send_data[1024],recv_data[1024];
        struct hostent *host;
        struct sockaddr_in server_addr;  

        host = gethostbyname(IP);

        if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            perror("Socket");
            exit(1);
        }

        server_addr.sin_family = AF_INET;     
        server_addr.sin_port = htons(port);   
        server_addr.sin_addr = *((struct in_addr *)host->h_addr);
        bzero(&(server_addr.sin_zero),8); 

        if (connect(sock, (struct sockaddr *)&server_addr,sizeof(struct sockaddr)) == -1) 
        {
            perror("Connect");
            exit(1);
        }

        while(1)
        {
        
          bytes_recieved=recv(sock,recv_data,1024,0);
          recv_data[bytes_recieved] = '\0';
 
          if (strcmp(recv_data , "q") == 0 || strcmp(recv_data , "Q") == 0)
          {
           close(sock);
           break;
          }

          else
           printf("\nRecieved data = %s " , recv_data);
           
           printf("\nSEND (q or Q to quit) : ");
           scanf("%s",send_data);
           
          if (strcmp(send_data , "q") != 0 && strcmp(send_data , "Q") != 0)
           send(sock,send_data,strlen(send_data), 0); 

          else
          {
           send(sock,send_data,strlen(send_data), 0);   
           close(sock);
           break;
          }
        
        }   
	
	return 0;
}

int main()
{	
  printf("Hello!\nWelcome to Chat Application\n");
  int boolean=1;
  printf("Enter your name: ");
  scanf("%s",self);
  printf("Enter your port: ");
  scanf("%d",&self_port);
  printf("Enter the name of your friend: ");
  scanf("%s",Name);
  printf("Enter the IP address of your friend: ");
  scanf("%s",IP);
  printf("Enter the port: ");
  scanf("%d",&port);
  // pthread_t t;
  // pthread_attr_t attr;
  // pthread_attr_init(&attr);
  // pthread_create(&t,&attr,client,NULL);
  server();
  
  
  // while(boolean)
  // {
    // printf("1. Add friend\n2. Leave chat room\n");
    // int i;
    // scanf("%d",&i);
    // switch(i)
    // {
    //   case 1:
    //   {
        
    //     printf("Enter the name of your friend: ");
    //     scanf("%s",Name);
    //     printf("Enter the IP address of your friend: ");
    //     scanf("%s",IP);
    //     printf("Enter the port: ");
    //     scanf("%d",&port);
    //     client(Name,IP,port);
    //     break;
    //   }

    //   case 2:
    //   boolean=0;
    //   break;

    //   default:
    //   boolean=0;
    //   //Do Nothing
    // }
  //}
  printf("Bye!\nSee you soon\n");
  
	return 0;
}














