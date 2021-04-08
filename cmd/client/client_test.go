package main

import (
	"github.com/fanjindong/gim/proto"
	"github.com/fanjindong/gim/proto/msg"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"net"
	"os"
	"testing"
	"time"
)

var (
	conn  net.Conn
	codec *proto.MsgProtocol
)

func TestMain(m *testing.M) {
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	codec = &proto.MsgProtocol{}

	code := m.Run()
	os.Exit(code)
}

type gConn struct {
	gnet.Conn
	ctx interface{}
}

func (g *gConn) Context() (ctx interface{}) {
	return g.ctx
}

func (g *gConn) SetContext(ctx interface{}) {
	g.ctx = ctx
}

func TestPing(t *testing.T) {
	req := msg.PingReq{
		Uid: "1",
	}
	data, _ := protobuf.Marshal(&req)

	t.Log(data, len(data))

	msg := proto.MsgProtocol{
		Version:    proto.DefaultProtocolVersion,
		Type:       proto.Ping,
		DataLength: uint16(len(data)),
		Data:       data,
	}
	body, err := codec.Encode(nil, msg.Bytes())
	if err != nil {
		panic(err)
	}
	t.Logf("body: %+v, length: %d", body, len(body))
	conn.Write(body)
	time.Sleep(1 * time.Second)
	pong := make([]byte, 0, len(body))
	conn.Read(pong)
	resp, err := proto.NewMsgProtocolByBytes(pong)
	t.Logf("pong: %+v, %+v, %v\n", pong, resp, err)
}
