package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	pb "github.com/sjansen/mecha/spikes/grpc/translate"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	address = "127.0.0.1:50051"
)

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type server struct{}

func (s *server) Translate(ctx context.Context, in *pb.Original) (*pb.Translation, error) {
	sb := strings.Builder{}
	for _, ch := range in.Msg {
		if ch >= 'a' && ch <= 'm' {
			ch += 13
		} else if ch >= 'n' && ch <= 'z' {
			ch -= 13
		}
		sb.WriteRune(ch)
	}
	return &pb.Translation{Msg: sb.String()}, nil
}

func startServer() {
	log.Println("Server starting...")

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTranslatorServer(s, &server{})
	reflection.Register(s)
	go func() {
		time.Sleep(5 * time.Second)
		s.GracefulStop()
	}()

	log.Println("Server started.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Println("Server stopped.")
}

func startClient(msg string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTranslatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Translate(ctx, &pb.Original{Msg: msg})
	if err != nil {
		log.Fatalf("could not translate: %v", err)
	}
	fmt.Printf("%s\t→ %s\n", msg, r.Msg)
}

func startChildren() {
	server := exec.Command(os.Args[0], "--as-server")
	server.Stdin = nil
	server.Stdout = os.Stdout
	server.Stderr = os.Stderr
	if err := server.Start(); err != nil {
		die(err)
	}

	time.Sleep(time.Second)
	client := exec.Command(os.Args[0], "--as-client")
	client.Stdin = nil
	client.Stdout = os.Stdout
	client.Stderr = os.Stderr
	if err := client.Start(); err != nil {
		die(err)
	}

	if err := client.Wait(); err != nil {
		die(err)
	}
	if err := server.Wait(); err != nil {
		die(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		startChildren()
	} else if os.Args[1] == "--as-server" {
		startServer()
	} else {
		for _, msg := range []string{
			"sbb", "one", "onm", "dhk", "dhhk", "pbetr", "tenhyg", "tnecyl", "jnyqb", "serq",
		} {
			startClient(msg)
		}
	}
}
