package packet

import (
	"fmt"

	"sutext.github.io/entry/types"
)

type ConnectPacket struct {
	UserID      string
	Platform    types.Platform
	AccessToken string
}

func Connect(userID string, platfrom types.Platform, accessToken string) *ConnectPacket {
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
	buffer := NewBuffer([]byte{})
	buffer.WriteString(p.UserID)
	buffer.WriteUInt8(byte(p.Platform))
	buffer.WriteString(p.AccessToken)
	return buffer.Bytes()
}
func (p *ConnectPacket) decode(data []byte) error {
	buffer := NewBuffer(data)
	userID, err := buffer.ReadString()
	if err != nil {
		return err
	}
	platform, err := buffer.ReadUInt8()
	if err != nil {
		return err
	}
	accessToken, err := buffer.ReadString()
	if err != nil {
		return err
	}
	p.Platform = types.Platform(platform)
	p.UserID = userID
	p.AccessToken = accessToken
	return nil
}

type ConnackCode uint16

const (
	// Connection Accepted
	ConnectionAccepted ConnackCode = 0
	// Unacceptable protocol version
	UnacceptableProtocolVersion ConnackCode = 1
	// Identifier rejected
	IdentifierRejected ConnackCode = 2
	// Server unavailable
	ServerUnavailable ConnackCode = 3
	// Bad user name or password
	BadUserNameOrPassword ConnackCode = 4
	// Not authorized
	NotAuthorized ConnackCode = 5
)

type ConnackPacket struct {
	Code ConnackCode
}

func Connack(code ConnackCode) *ConnackPacket {
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
	buffer := NewBuffer([]byte{})
	buffer.WriteUInt16(uint16(p.Code))
	return buffer.Bytes()
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
	buffer := NewBuffer(data)
	code, err := buffer.ReadUInt16()
	if err != nil {
		return err
	}
	p.Code = ConnackCode(code)
	return nil
}
