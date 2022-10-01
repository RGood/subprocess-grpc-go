package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/RGood/subprocesses-go/internal/dialer"
	"github.com/RGood/subprocesses-go/internal/protos/echo"
	"github.com/RGood/subprocesses-go/internal/protos/runtime"
	echoImpl "github.com/RGood/subprocesses-go/pkg/echo"
	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		panic(errors.New("Subprocess requires socket and runtime address"))
	}

	id := args[0]

	socketAddr := fmt.Sprintf("/tmp/%s.sock", id)
	runtimeID := args[1]

	os.Remove(socketAddr)
	//lis, err := net.Listen("tcp", port)
	serverAddr, err := net.ResolveUnixAddr("unix", socketAddr)
	if err != nil {
		log.Fatal("failed to resolve unix addr")
	}

	lis, err := net.ListenUnix("unix", serverAddr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	echo.RegisterEchoServiceServer(s, &echoImpl.EchoServiceImpl{})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.Serve(lis)
		wg.Done()
	}()

	runtimeConn, _ := dialer.Dial(runtimeID)
	runtimeClient := runtime.NewRuntimeClient(runtimeConn)
	runtimeClient.Ready(context.Background(), &runtime.ReadyMessage{
		SubprocessId: id,
	})

	wg.Wait()
}
