package main

import (
	//"context"

	"net"
	"os"

	"github.com/urfave/cli"
	pb "github.com/yulefox/lamp/kits/cli/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SayHello(context.Context, *pb.Hello_Req) (*pb.Hello_Res, error) {
	return nil, nil
}

func main() {
	app := &cli.App{
		Name: "cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: ":12345",
				Usage: "listening address host:port",
			},
		},
		Action: func(c *cli.Context) error {
			lis, err := net.Listen("tcp", c.String("addr"))
			if err != nil {
				panic(err)
			}
			s := grpc.NewServer()
			srv := &server{}
			pb.RegisterGreeterServer(s, srv)
			return s.Serve(lis)
		},
	}
	app.Run(os.Args)
}
