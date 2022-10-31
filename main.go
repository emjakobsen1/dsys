package main

import (
	"log"
	"net"
	"sync"

	"example.com/proto"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Connection struct {
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
	proto.UnimplementedBroadcastServer
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	log.Println("CreateStream: ", conn)
	return <-conn.error
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)

				log.Println("Broadcasting message to: ", conn.stream, msg)

				if err != nil {
					log.Fatalf("Error with Stream: %v - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	log.Println("Close: ")
	return &proto.Close{}, nil

}

func main() {
	var connections []*Connection
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	log.Println("Starting Chitty-Chat at port :8080")
	server := &Server{connections, proto.UnimplementedBroadcastServer{}}
	proto.RegisterBroadcastServer(grpcServer, server)
	log.Println("RegisterBroadcastServer: ", grpcServer, server)
	grpcServer.Serve(listener)
}
