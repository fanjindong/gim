package proto

//
// import (
// 	"encoding/binary"
// 	"github.com/fanjindong/gim/errors"
// 	"github.com/fanjindong/gim/proto/msg"
// 	"github.com/panjf2000/gnet"
// )
//
// type MsgCodec struct{}
//
// func (MsgCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
// 	panic("implement me")
// }
//
// func (MsgCodec) Decode(c gnet.Conn) ([]byte, error) {
// 	size, buf := c.ReadN(1)
// 	msgTypeInt := int32(binary.BigEndian.Uint16(buf))
// 	if _, ok := msg.MsgType_name[msgTypeInt]; !ok {
// 		return nil, errors.UnknownMsgType
// 	}
// 	msgType := msg.MsgType(msgTypeInt)
//
// 	c.ShiftN(size)
// 	size, buf = c.ReadN(2)
// 	bodyLength := int(binary.BigEndian.Uint16(buf))
// 	c.ShiftN(size)
//
// 	var body []byte
// 	if bodyLength > 0 {
// 		size, body = c.ReadN(bodyLength)
// 	}
//
// 	msg := NewMsg(msgType, body)
// 	return
// }
//
// type Msg struct {
// 	Type msg.MsgType
// 	Data []byte
// }
//
// func NewMsg(t msg.MsgType, data []byte) Msg {
// 	return Msg{Type: t, Data: data}
// }
//
// func (m *Msg) Bytes() []byte {
// 	jsoniters
// }
//
//
