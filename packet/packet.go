package packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"sutext.github.io/entry/buffer"
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
	PING    PacketType = 3 // PING packet
	PONG    PacketType = 4 // PONG packet
	DATA    PacketType = 5 // DATA packet
	CLOSE   PacketType = 6 // CLOSE packet
)

func (t PacketType) String() string {
	switch t {
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case PING:
		return "PING"
	case PONG:
		return "PONG"
	case DATA:
		return "DATA"
	case CLOSE:
		return "CLOSE"
	default:
		return "UNKNOWN"
	}
}

type Packet interface {
	fmt.Stringer
	buffer.WriteTo
	buffer.ReadFrom
	Type() PacketType
	Equal(Packet) bool
}

type incomming struct {
	packetType PacketType
	data       []byte
}

func (i incomming) decode() (Packet, error) {
	buf := buffer.New(i.data)
	switch i.packetType {
	case CONNECT:
		conn := &ConnectPacket{}
		err := conn.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return conn, nil
	case CONNACK:
		connack := &ConnackPacket{}
		err := connack.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return connack, nil
	case DATA:
		data := &DataPacket{}
		err := data.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		return data, nil
	case PING:
		return Ping(), nil
	case PONG:
		return Pong(), nil
	case CLOSE:
		close := &ClosePacket{}
		err := close.ReadFrom(buf)
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
	buf := buffer.New()
	err := p.WriteTo(buf)
	if err != nil {
		return err
	}
	// write header
	length := uint32(buf.Len())
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
	_, err = w.Write(header)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

type packet struct {
	packetType PacketType
	data       []byte
}

func (p packet) Type() PacketType {
	return p.packetType
}

func (p packet) Equal(other packet) bool {
	if p.packetType != other.Type() {
		return false
	}
	return reflect.DeepEqual(p.data, other.data)
}
func (p packet) String() string {
	return p.String()
}
