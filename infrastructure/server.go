package infrastructure

import (
	"bufio"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func Init(host string, port string, bind func(server *grpc.Server)) {
	create(host, port, bind)
}

func create(host string, port string, bind func(server *grpc.Server)) {

	server := grpc.NewServer()

	bind(server)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("Failed to listen on port 8080: %s", err)
	}

	if err = server.Serve(lis); err != nil {
		log.Println(err)
	} else {
		log.Printf("Starting server at %s", lis.Addr().Network())
	}

	reader := bufio.NewReader(os.Stdin)
	if text, _ := reader.ReadString('\n'); text == "stop" {
		server.GracefulStop()
	}
}
