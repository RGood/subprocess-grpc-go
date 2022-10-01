package dialer

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

func Dial(processID string) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("/tmp/%s.sock", processID)
	return grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDialer(dialer))
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}

	return net.DialUnix("unix", nil, unixAddr)
}
