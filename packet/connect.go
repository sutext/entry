package packet

import (
	"fmt"
)

const (
	ConnectFlagIdentity uint8 = 0x01
)

type Identity struct {
	Token    string
	UserID   string
	ClientID string
}
type ConnectPacket struct {
	Identity *Identity
	flag     uint8
}

func Connect(identity *Identity) *ConnectPacket {
	var flag uint8 = 0
	if identity != nil {
		flag |= ConnectFlagIdentity
	}
	return &ConnectPacket{Identity: identity, flag: flag}
}
func (p *ConnectPacket) String() string {
	return fmt.Sprintf("CONNECT(uid=%s, cid=%s, token=%s)", p.Identity.UserID, p.Identity.ClientID, p.Identity.Token)
}
func (p *ConnectPacket) Type() PacketType {
	return CONNECT
}
func (p *ConnectPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	if other.Type() != CONNECT {
		return false
	}
	otherP := other.(*ConnectPacket)
	return p.Identity.Token == otherP.Identity.Token &&
		p.Identity.UserID == otherP.Identity.UserID &&
		p.Identity.ClientID == otherP.Identity.ClientID &&
		p.flag == otherP.flag
}
func (p *ConnectPacket) encode() []byte {
	buffer := newBuffer([]byte{})
	buffer.writeUInt8(p.flag)
	if p.flag&ConnectFlagIdentity != 0 {
		buffer.writeString(p.Identity.Token)
		buffer.writeString(p.Identity.UserID)
		buffer.writeString(p.Identity.ClientID)
	}
	return buffer.bytes()
}
func (p *ConnectPacket) decode(data []byte) error {
	buffer := newBuffer(data)
	flag, err := buffer.readUInt8()
	if err != nil {
		return err
	}
	p.flag = flag
	if flag&ConnectFlagIdentity != 0 {
		token, err := buffer.readString()
		if err != nil {
			return err
		}
		userID, err := buffer.readString()
		if err != nil {
			return err
		}
		clientID, err := buffer.readString()
		if err != nil {
			return err
		}
		p.Identity = &Identity{
			Token:    token,
			UserID:   userID,
			ClientID: clientID,
		}
	}
	return nil
}

type ConnectCode uint16

const (
	// Connection Accepted
	ConnectionAccepted ConnectCode = 0
	// Connection Refused, unacceptable protocol version
	AlreadyConnected ConnectCode = 1
	// Identifier rejected
	IdentifierRejected ConnectCode = 2
	// Server unavailable
	ServerUnavailable ConnectCode = 3
	// Bad user name or password
	BadUserNameOrPassword ConnectCode = 4
	// Not authorized
	NotAuthorized ConnectCode = 5
)

type ConnackPacket struct {
	Code ConnectCode
}

func Connack(code ConnectCode) *ConnackPacket {
	return &ConnackPacket{
		Code: code,
	}
}
func (p *ConnackPacket) String() string {
	return fmt.Sprintf("CONNACK(%d)", p.Code)
}
func (P *ConnackPacket) Type() PacketType {
	return CONNACK
}
func (p *ConnackPacket) encode() []byte {
	buffer := newBuffer([]byte{})
	buffer.writeUInt16(uint16(p.Code))
	return buffer.bytes()
}
func (p *ConnackPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	if other.Type() != CONNACK {
		return false
	}
	otherP := other.(*ConnackPacket)
	return p.Code == otherP.Code
}
func (p *ConnackPacket) decode(data []byte) error {
	buffer := newBuffer(data)
	code, err := buffer.readUInt16()
	if err != nil {
		return err
	}
	p.Code = ConnectCode(code)
	return nil
}
