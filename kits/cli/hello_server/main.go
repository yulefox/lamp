package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"

	pb "github.com/yulefox/lamp/kits/cli/proto"
	"golang.org/x/net/context"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *pb.Hello_Req) (*pb.Hello_Res, error) {
	return &pb.Hello_Res{Message: "hello" + in.Name}, nil
}

func (s *server) Tunnel(stream pb.GameService_TunnelServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(in); err != nil {
			panic(err)
		}
		fmt.Println(in)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":10101")

	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	pb.RegisterGameServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
