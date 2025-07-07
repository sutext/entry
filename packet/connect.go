package packet

import (
	"fmt"

	"sutext.github.io/entry/buffer"
)

const (
	ConnectFlagIdentity uint8 = 0x01
)

type Identity struct {
	UserID    string
	ClientID  string
	AuthToken string
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
	return fmt.Sprintf("CONNECT(uid=%s, cid=%s, token=%s)", p.Identity.UserID, p.Identity.ClientID, p.Identity.AuthToken)
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
	return p.Identity.AuthToken == otherP.Identity.AuthToken &&
		p.Identity.UserID == otherP.Identity.UserID &&
		p.Identity.ClientID == otherP.Identity.ClientID &&
		p.flag == otherP.flag
}
func (p *ConnectPacket) WriteTo(buffer *buffer.Buffer) error {
	buffer.WriteUInt8(p.flag)
	if p.flag&ConnectFlagIdentity != 0 {
		buffer.WriteString(p.Identity.AuthToken)
		buffer.WriteString(p.Identity.UserID)
		buffer.WriteString(p.Identity.ClientID)
	}
	return nil
}
func (p *ConnectPacket) ReadFrom(buffer *buffer.Buffer) error {
	flag, err := buffer.ReadUInt8()
	if err != nil {
		return err
	}
	p.flag = flag
	if flag&ConnectFlagIdentity != 0 {
		token, err := buffer.ReadString()
		if err != nil {
			return err
		}
		userID, err := buffer.ReadString()
		if err != nil {
			return err
		}
		clientID, err := buffer.ReadString()
		if err != nil {
			return err
		}
		p.Identity = &Identity{
			AuthToken: token,
			UserID:    userID,
			ClientID:  clientID,
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
func (p *ConnackPacket) WriteTo(buffer *buffer.Buffer) error {
	buffer.WriteUInt16(uint16(p.Code))
	return nil
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
func (p *ConnackPacket) ReadFrom(buffer *buffer.Buffer) error {
	code, err := buffer.ReadUInt16()
	if err != nil {
		return err
	}
	p.Code = ConnectCode(code)
	return nil
}
