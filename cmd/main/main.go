package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/RGood/subprocesses-go/internal/dialer"
	"github.com/RGood/subprocesses-go/internal/helpers"
	"github.com/RGood/subprocesses-go/internal/protos/echo"
	runtimeGRPC "github.com/RGood/subprocesses-go/internal/protos/runtime"
	"github.com/RGood/subprocesses-go/pkg/runtime"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func startService(r *runtime.RuntimeService, subprocessID, runtimeID string, commandParts []string) *exec.Cmd {
	rootCommand := commandParts[0]
	processArgs := commandParts[1:]
	processArgs = append(processArgs, subprocessID, runtimeID)

	r.AddService(subprocessID)
	// Create command to run subprocess
	cmd := exec.Command(rootCommand, processArgs...)

	// Start subprocess (non-blocking)
	cmd.Start()

	// We should add a timeout to this
	r.WaitForService(subprocessID)

	return cmd
}

func main() {
	// Create subprocess ID
	subprocessID := uuid.New().String()

	// Create runtime ID and address
	runtimeID := uuid.New().String()
	runtimeAddress := helpers.CreateSocketAddress(runtimeID)
	r := runtime.NewRuntime()

	// Load command to run subprocess
	commandParts := strings.Split(os.Getenv("PROCESS_CMD"), " ")

	// Clear runtime socket
	os.Remove(runtimeAddress)

	// Resolve runtime socket
	runtimeAddr, err := net.ResolveUnixAddr("unix", runtimeAddress)
	if err != nil {
		log.Fatal("fialed to resolve unix addr")
	}

	// Listen on runtime socket
	lis, err := net.ListenUnix("unix", runtimeAddr)
	if err != nil {
		panic(err)
	}

	// Instantiate runtime grpc server
	s := grpc.NewServer()

	// Register runtime server on the grpc server
	runtimeGRPC.RegisterRuntimeServer(s, r)

	// Listen with the grpc server in the background
	go s.Serve(lis)

	// Start the remote service on the given subprocess ID
	//   wait for it to report ready, and return the command
	//   for process control
	start := time.Now()
	cmd := startService(r, subprocessID, runtimeID, commandParts)
	fmt.Printf("Service started in %v\n", time.Since(start))

	// Connect to remote process
	conn, err := dialer.Dial(subprocessID)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Instantiate client
	c := echo.NewEchoServiceClient(conn)

	// Call the client N times and wait for all responses to resolve
	wg := sync.WaitGroup{}
	requestStart := time.Now()
	reqCount := 100
	for i := 0; i < reqCount; i++ {
		wg.Add(1)
		go func(idx int) {
			_, err := c.Echo(context.Background(), &echo.Message{Text: fmt.Sprintf("foo: %d", idx)})
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Avg request time: %v\n", time.Since(requestStart)/time.Duration(reqCount))

	// Terminate remote process and wait for it to conclude
	cmd.Process.Kill()
	cmd.Wait()
}
