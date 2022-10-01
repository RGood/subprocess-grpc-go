package dialer

import (
	"net"
	"time"

	"github.com/RGood/subprocesses-go/internal/helpers"
	"google.golang.org/grpc"
)

func Dial(processID string) (*grpc.ClientConn, error) {
	addr := helpers.CreateSocketAddress(processID)
	return grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDialer(dialer))
}

func dialer(addr string, timeout time.Duration) (net.Conn, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}

	return net.DialUnix("unix", nil, unixAddr)
}
