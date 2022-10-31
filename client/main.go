package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"example.com/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamerror error
	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})
	//log.Println(user)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("error reading message: %v", err)
				break
			}
			user.T++
			//max(t,t') +1
			if user.GetT() > msg.GetTimestamp() {
				user.T = user.GetT() + 1
			} else {
				user.T = msg.GetTimestamp() + 1
			}

			fmt.Printf("%s #%v (T = %v)%s\n", msg.GetUserName(), msg.GetId(), user.GetT(), msg.GetContent())

		}
	}(stream)

	return streamerror
}

func main() {

	done := make(chan int)

	name := flag.String("name", "Anonymous user", "type -name to set the name of the user")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(100)
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	client = proto.NewBroadcastClient(conn)
	user := &proto.User{
		Id:   fmt.Sprintf("%d", id),
		Name: *name,
		T:    0, //on initialisation
	}

	connect(user)
	wait.Add(1)

	go func() {
		user.T++ //increment on event
		client.BroadcastMessage(context.Background(), &proto.Message{
			Id:        user.Id,
			Content:   " joined Chitty-Chat server.",
			Timestamp: user.GetT(),
			UserName:  user.GetName(),
		})
		defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			user.T++ //increment on event
			msg := &proto.Message{
				Id:        user.Id,
				Content:   ":" + scanner.Text(),
				Timestamp: user.GetT(),
				UserName:  user.GetName(),
			}
			if msg.GetContent() == ":-close" {
				client.BroadcastMessage(context.Background(), &proto.Message{
					Id:        user.Id,
					Content:   " left the server.",
					Timestamp: user.GetT(),
					UserName:  user.GetName(),
				})
				close(done)
			}
			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				log.Printf("Error Sending Message: %v", err)
				break
			}
		}

	}()

	go func() {
		wait.Wait()

		close(done)

	}()

	<-done
}
