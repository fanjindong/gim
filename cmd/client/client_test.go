package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/fanjindong/gim/proto"
	"github.com/fanjindong/gim/proto/msg"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"io"
	"net"
	"os"
	"testing"
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

	reqMsg := proto.MsgProtocol{
		Version:    proto.DefaultProtocolVersion,
		Type:       proto.Ping,
		DataLength: uint16(len(data)),
		Data:       data,
	}
	body, err := codec.Encode(nil, reqMsg.Bytes())
	if err != nil {
		panic(err)
	}
	t.Logf("body: %+v, length: %d", body, len(body))
	conn.Write(body)
	// time.Sleep(1 * time.Second)
	respMsg := readByConn(conn)
	pong := &msg.PingReply{}
	protobuf.Unmarshal(respMsg.Data, pong)
	// t.Logf("resp: %+v, \npong: %+v\n", respMsg, pong)
	t.Log("code:", pong.Code, "msg:", pong.Msg)
}

func readByConn(conn io.Reader) *proto.MsgProtocol {
	mp := &proto.MsgProtocol{}
	header := make([]byte, proto.DefaultHeadLength)
	conn.Read(header)
	fmt.Println(header)
	buffer := bytes.NewBuffer(header)

	_ = binary.Read(buffer, binary.BigEndian, &mp.Version)
	_ = binary.Read(buffer, binary.BigEndian, &mp.Type)
	_ = binary.Read(buffer, binary.BigEndian, &mp.DataLength)

	if mp.DataLength > 0 {
		data := make([]byte, mp.DataLength)
		conn.Read(data)
		mp.Data = data
	}
	return mp
}
