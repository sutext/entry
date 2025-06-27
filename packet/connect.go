package packet

import (
	"fmt"

	"sutext.github.io/entry/code"
)

type ConnectPacket struct {
	UserID      string
	Platform    code.Platform
	AccessToken string
}

func Connect(userID string, platfrom code.Platform, accessToken string) *ConnectPacket {
	return &ConnectPacket{
		UserID:      userID,
		Platform:    platfrom,
		AccessToken: accessToken,
	}
}

func (p *ConnectPacket) String() string {
	return fmt.Sprintf("CONNECT(uid=%s,platform=%d, token=%s)", p.UserID, p.Platform, p.AccessToken)
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
	return p.Platform == otherP.Platform && p.UserID == otherP.UserID && p.AccessToken == otherP.AccessToken
}
func (p *ConnectPacket) encode() []byte {
	buffer := newBuffer([]byte{})
	buffer.writeString(p.UserID)
	buffer.writeUInt8(byte(p.Platform))
	buffer.writeString(p.AccessToken)
	return buffer.bytes()
}
func (p *ConnectPacket) decode(data []byte) error {
	buffer := newBuffer(data)
	userID, err := buffer.readString()
	if err != nil {
		return err
	}
	platform, err := buffer.readUInt8()
	if err != nil {
		return err
	}
	accessToken, err := buffer.readString()
	if err != nil {
		return err
	}
	p.Platform = code.Platform(platform)
	p.UserID = userID
	p.AccessToken = accessToken
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
