package service

import (
	"github.com/fanjindong/gim/proto"
	"github.com/fanjindong/gim/proto/msg"
	rs "github.com/fanjindong/gim/resource"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
)

type Server struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func NewServer(addr string, multicore bool, async bool, codec gnet.ICodec, workerPool *goroutine.Pool) *Server {
	return &Server{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: workerPool}
}

func (cs *Server) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cs *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	m, _ := proto.NewMsgProtocolByBytes(frame)
	req := &msg.PingReq{}
	defer rs.Logger.Infof("React: %+v, %+v, %+v\n", frame, m, req)
	if err := protobuf.Unmarshal(m.Data, req); err != nil {
		rs.Logger.Errorln("React.Unmarshal err: ", err.Error())
		return
	}

	reply := proto.NewMsgProtocol(proto.Ping)
	if err := reply.SetData(&msg.PingReply{
		Code: 0,
		Msg:  "欢迎 " + req.Uid + " 用户",
	}); err != nil {
		rs.Logger.Errorln("React.reply.SetData err: ", err.Error())
		return
	}

	// if cs.async {
	// 	_ = cs.workerPool.Submit(func() {
	// 		c.AsyncWrite(reply.Bytes())
	// 	})
	// 	return
	// }
	out = reply.Bytes()
	return
}
