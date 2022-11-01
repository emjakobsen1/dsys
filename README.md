### ChittyChat
#### About
*ChittyChat* is a console-based chat service that lets multiple clients communicate through a server. It is implemented in Go and uses gRPC for passing messages between participants. It computes and logs Lamport timestamps throughout the lifetime of clients.

#### Running the application
Clone the project to a local repository. From the root repository, open a terminal and run the server with 

```
 go run main.go
```
In another terminal, Change directory to the client repository, and start initialize client with 
``` 
go run main.go -name <yourName>
```
The flag ```-name``` followed by a custom name is optional. 

You can now chat via the command line. To end a client, type ```-close```. 
