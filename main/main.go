package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	rabbit "github.com/cs-sea/rabbit/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	amqppool "github.com/cs-sea/rabbit"
	_ "github.com/cs-sea/rabbit/compile"
)

type APIServer struct {
	p amqppool.Producer
}

func NewAPIServer(p amqppool.Producer) *APIServer {
	return &APIServer{p: p}
}

func (a *APIServer) Push(ctx context.Context, req *rabbit.PushMessageRequest) (*rabbit.EmptyResponse, error) {
	fmt.Printf("%+v", req)
	err := a.p.Publish(ctx, &amqppool.Exchange{
		Name: req.Exchange.Name,
		Kind: req.Exchange.Kind.String(),
	}, &amqppool.Message{
		Body:       req.Message,
		RoutingKey: req.Key,
	})

	fmt.Println(err)
	return &rabbit.EmptyResponse{}, nil
}

var (
	serverDefaultStackSize = 4 << 10
)

func main() {

	logger := logrus.NewEntry(logrus.New())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8888))
	if err != nil {
		panic(err)
	}

	logger.Println("start rpc")

	gs := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 10 * time.Minute,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
				grpc_ctxtags.WithFieldExtractor(func(fullMethod string, req interface{}) map[string]interface{} {
					fields := map[string]interface{}{"request_id": xid.New().String()}
					return fields
				}),
			),
			grpc_logrus.UnaryServerInterceptor(logger),
			grpc_logrus.PayloadUnaryServerInterceptor(
				logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
					return true
				},
			),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
				err, ok := p.(error)
				if !ok {
					err = fmt.Errorf("%v", p)
				}
				stack := make([]byte, serverDefaultStackSize)
				length := runtime.Stack(stack, false)
				logger.WithError(err).Errorf("recovered from panic: %v [stack]: %s", err, stack[:length])
				return status.Errorf(codes.Internal, "Unexcepted internal server error")
			})),
		),
	)
	rabbitService := amqppool.NewRabbit(&amqppool.Options{})

	producer := amqppool.NewProducerService(rabbitService)
	apiServer := NewAPIServer(producer)
	rabbit.RegisterRabbitServiceServer(gs, apiServer)

	go gs.Serve(lis)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("exit")
	gs.GracefulStop()
}
