package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	rs "github.com/fanjindong/gim/resource"
	protobuf "github.com/golang/protobuf/proto"
	"log"

	"github.com/panjf2000/gnet"
)

// MsgProtocol : custom protocol
// custom protocol header contains Version, ActionType and DataLength fields
// its payload is Data field
type MsgProtocol struct {
	Version    uint8
	Type       uint8
	DataLength uint16
	Data       []byte
}

func NewMsgProtocol(msgType uint8) *MsgProtocol {
	return &MsgProtocol{Type: msgType, Version: DefaultProtocolVersion}
}

func NewMsgProtocolByBytes(data []byte) (*MsgProtocol, error) {
	msg := &MsgProtocol{}
	buffer := bytes.NewBuffer(data[:DefaultHeadLength+1])

	_ = binary.Read(buffer, binary.BigEndian, &msg.Version)
	_ = binary.Read(buffer, binary.BigEndian, &msg.Type)
	_ = binary.Read(buffer, binary.BigEndian, &msg.DataLength)

	if msg.DataLength > 0 {
		msg.Data = data[DefaultHeadLength:]
	}
	return msg, nil
}

func (cc *MsgProtocol) SetData(message protobuf.Message) error {
	data, err := protobuf.Marshal(message)
	if err != nil {
		return err
	}
	cc.Data = data
	cc.DataLength = uint16(len(data))
	return nil
}

func (cc *MsgProtocol) Bytes() []byte {
	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.BigEndian, cc.Version); err != nil {
		return nil
	}
	if err := binary.Write(buffer, binary.BigEndian, cc.Type); err != nil {
		return nil
	}
	if err := binary.Write(buffer, binary.BigEndian, cc.DataLength); err != nil {
		return nil
	}
	if cc.DataLength > 0 {
		if err := binary.Write(buffer, binary.BigEndian, cc.Data); err != nil {
			return nil
		}
	}
	return buffer.Bytes()
}

// Encode ...
func (cc *MsgProtocol) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	result := make([]byte, 0)
	msg, err := NewMsgProtocolByBytes(buf)

	defer rs.Logger.Infoln("codec.Encode: ", msg)

	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(result)

	if err := binary.Write(buffer, binary.BigEndian, msg.Version); err != nil {
		s := fmt.Sprintf("Pack version error , %v", err)
		return nil, errors.New(s)
	}
	if err := binary.Write(buffer, binary.BigEndian, msg.Type); err != nil {
		s := fmt.Sprintf("Pack type error , %v", err)
		return nil, errors.New(s)
	}
	if err := binary.Write(buffer, binary.BigEndian, msg.DataLength); err != nil {
		s := fmt.Sprintf("Pack datalength error , %v", err)
		return nil, errors.New(s)
	}
	if msg.DataLength > 0 {
		if err := binary.Write(buffer, binary.BigEndian, msg.Data); err != nil {
			s := fmt.Sprintf("Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	return buffer.Bytes(), nil
}

// Decode ...
func (cc *MsgProtocol) Decode(c gnet.Conn) ([]byte, error) {
	msg := MsgProtocol{}

	defer rs.Logger.Infoln("codec.Decode: ", &msg)

	// parse header
	headerLen := DefaultHeadLength // uint8+uint8+uint16
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		_ = binary.Read(byteBuffer, binary.BigEndian, &msg.Version)
		_ = binary.Read(byteBuffer, binary.BigEndian, &msg.Type)
		_ = binary.Read(byteBuffer, binary.BigEndian, &msg.DataLength)

		// to check the protocol version and msgType,
		// reset buffer if the version or msgType is not correct
		if msg.Version != DefaultProtocolVersion || isCorrectAction(msg.Type) == false {
			c.ResetBuffer()
			log.Println("not normal protocol:", msg.Version, DefaultProtocolVersion, msg.Type, msg.DataLength)
			return nil, errors.New("not normal protocol")
		}
		// parse payload
		dataLen := int(msg.DataLength) // max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen
		if dataSize, data := c.ReadN(protocolLen); dataSize == protocolLen {
			c.ShiftN(protocolLen)
			// log.Println("parse success:", data, dataSize)
			msg.Data = data[headerLen:]
			// return the payload of the data
			return msg.Bytes(), nil
		}
		// log.Println("not enough payload data:", dataLen, protocolLen, dataSize)
		return nil, errors.New("not enough payload data")

	}
	// log.Println("not enough header data:", size)
	return nil, errors.New("not enough header data")
}

// default custom protocol const
const (
	DefaultHeadLength = 4

	DefaultProtocolVersion = 0x01 // test protocol version

	Ping = 0x00
	Send = 0x01
	Push = 0x02
)

func isCorrectAction(msgType uint8) bool {
	switch msgType {
	case Ping, Send, Push:
		return true
	default:
		return false
	}
}
