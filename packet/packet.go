package packet

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MIN_LEN uint32 = 0
	MID_LEN uint32 = 0xfff
	MAX_LEN uint32 = 0xfffffff
)

type PacketType uint8

const (
	CONNECT PacketType = 1 // CONNECT packet
	CONNACK PacketType = 2 // CONNACK packet
	DATA    PacketType = 3 // DATA packet
	DATAACK PacketType = 4 // DATAACK packet
	PING    PacketType = 5 // PING packet
	PONG    PacketType = 6 // PONG packet
	CLOSE   PacketType = 7 // CLOSE packet
)

func (t PacketType) String() string {
	switch t {
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case DATA:
		return "DATA"
	case DATAACK:
		return "DATAACK"
	case PING:
		return "PING"
	case PONG:
		return "PONG"
	case CLOSE:
		return "CLOSE"
	default:
		return "UNKNOWN"
	}
}

type Packet interface {
	fmt.Stringer
	Type() PacketType
	Equal(Packet) bool
	encode() []byte
	decode(data []byte) error
}

type incomming struct {
	packetType PacketType
	data       []byte
}

func (i incomming) decode() (Packet, error) {
	switch i.packetType {
	case CONNECT:
		conn := &ConnectPacket{}
		err := conn.decode(i.data)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case CONNACK:
		connack := &ConnackPacket{}
		err := connack.decode(i.data)
		if err != nil {
			return nil, err
		}
		return connack, nil
	case DATA:
		data := &DataPacket{}
		err := data.decode(i.data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case DATAACK:
		dataack := &DataAckPacket{}
		err := dataack.decode(i.data)
		if err != nil {
			return nil, err
		}
		return dataack, nil
	case PING:
		return Ping(), nil
	case PONG:
		return Pong(), nil
	case CLOSE:
		close := &ClosePacket{}
		err := close.decode(i.data)
		if err != nil {
			return nil, err
		}
		return close, nil
	default:
		return nil, ErrUnkownPacketType
	}
}

func ReadPacket(r io.Reader) (Packet, error) {
	// read header
	header := make([]byte, 2)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return nil, err
	}
	packetType := PacketType(header[0] >> 5)
	//read length
	flag := header[0] & 0x10
	var length uint32
	if flag != 0 {
		bs := make([]byte, 2)
		io.ReadFull(r, bs)
		length = binary.BigEndian.Uint32([]byte{header[0] & 0x0f, header[1], bs[0], bs[1]})
	} else {
		length = binary.BigEndian.Uint32([]byte{0, 0, header[0] & 0x0f, header[1]})
	}
	// read data
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	imcoming := incomming{packetType: packetType, data: data}
	return imcoming.decode()
}
func WritePacket(w io.Writer, p Packet) error {
	data := p.encode()
	// write header
	length := uint32(len(data))
	if length > MAX_LEN {
		return ErrPacketSizeTooLarge
	}
	var header []byte
	if length > MID_LEN {
		header = make([]byte, 4)
		binary.BigEndian.PutUint32(header, length)
		header[0] = byte(p.Type()<<5) | 0x10 | header[0]
	} else {
		header = make([]byte, 2)
		binary.BigEndian.PutUint16(header, uint16(length))
		header[0] = byte(p.Type()<<5) | header[0]
	}
	// write data
	_, err := w.Write(header)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}
