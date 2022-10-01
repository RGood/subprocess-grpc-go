package echo

import (
	"context"
	"strings"

	"github.com/RGood/subprocesses-go/internal/protos/echo"
)

type EchoServiceImpl struct {
	echo.UnimplementedEchoServiceServer
}

func (es *EchoServiceImpl) Echo(ctx context.Context, msg *echo.Message) (*echo.Message, error) {
	msg.Text = strings.ToUpper(msg.Text)
	return msg, nil
}
