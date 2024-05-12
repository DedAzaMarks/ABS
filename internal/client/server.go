package client

import (
	"log"

	pb "github.com/DedAzaMarks/ABS/internal/proto"
)

const pong = "pong"

type Server struct {
	pb.UnimplementedClientServer
	DoPing chan struct{}
}

func (s *Server) Ping(_ *pb.Init, pingStream pb.Client_PingServer) error {
	for range s.DoPing {
		log.Print("ping client")
		err := pingStream.Send(&pb.Ping{})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewServer() *Server {
	return &Server{
		DoPing: make(chan struct{}),
	}
}
