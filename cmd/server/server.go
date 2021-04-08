package main

import (
	"flag"
	"fmt"
	"github.com/fanjindong/gim/proto"
	"github.com/fanjindong/gim/service"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
)

func main() {
	var port int
	var multicore bool

	// Example command: go run server.go --port 9000 --multicore=true
	flag.IntVar(&port, "port", 9000, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()
	addr := fmt.Sprintf("tcp://:%d", port)

	cs := service.NewServer(addr, multicore, false, &proto.MsgProtocol{}, goroutine.Default())
	err := gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(&proto.MsgProtocol{}))
	if err != nil {
		panic(err)
	}
}
