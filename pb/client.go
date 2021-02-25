package pb

import "google.golang.org/grpc"

func NewClient(target string) RabbitServiceClient {

	conn, err := grpc.Dial("localhost:8888")
	if err != nil {
		panic("ping service err")
	}

	return NewRabbitServiceClient(conn)
}
