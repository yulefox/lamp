package main

import (
	"log"

	"google.golang.org/grpc"

	"os"

	pb "github.com/yulefox/lamp/kits/cli/proto"
	"golang.org/x/net/context"
)

func main() {
	conn, err := grpc.Dial("localhost:10101", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.Hello_Req{Name: name})
	if err != nil {
		panic(err)
	}
	log.Println("reply: ", r.Message)
}
