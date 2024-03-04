package utils

import (
	"github.com/fatedier/frp/client"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func FrpTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}
