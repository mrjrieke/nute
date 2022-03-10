package server

import (
	"context"

	"google.golang.org/grpc"
	pb "nute/mashupsdk"
)

// server is used to implement server.MashupServer.
type MashupServer struct {
	pb.UnimplementedMashupServerServer
}

func (s *MashupServer) Shutdown(ctx context.Context, in *pb.MashupEmpty) (*pb.MashupEmpty, error) {
	os.Exit(-1)
    return &pb.MashupEmpty{}, nil
}

func InitServer() {
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

}