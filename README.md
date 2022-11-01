### ChittyChat
##### About
*ChittyChat* is a console-based chat service that lets multiple clients communicate through a server. It is implemented in Go and uses gRPC for passing messages between participants. It computes and logs Lamport timestamps throughout the lifetime of clients.

Clone the project to a local repository. From the root repository, run the server with 

```
 go run main.go
```
Change directory to the client repository, and start the clients with 
``` 
go run main.go -name <yourName>
```
The flag -name followed by your custom name is optional. 
