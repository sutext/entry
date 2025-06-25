package packet

import "fmt"

type CloseCode uint16

const (
	CloseNormal             CloseCode = 1000
	CloseGoingAway          CloseCode = 1001
	CloseProtocolError      CloseCode = 1002
	CloseUnsupported        CloseCode = 1003
	CloseNoStatus           CloseCode = 1005
	CloseAbnormalClosure    CloseCode = 1006
	CloseInvalidFrame       CloseCode = 1007
	ClosePolicyViolation    CloseCode = 1008
	CloseMessageTooBig      CloseCode = 1009
	CloseMandatoryExtension CloseCode = 1010
	CloseInternalError      CloseCode = 1011
	CloseServiceRestart     CloseCode = 1012
	CloseTryAgainLater      CloseCode = 1013
	CloseTLSHandshake       CloseCode = 1015
)

type ClosePacket struct {
	Code CloseCode
}

func Close(code CloseCode) *ClosePacket {
	return &ClosePacket{Code: code}
}
func (p *ClosePacket) String() string {
	return fmt.Sprintf("CLOSE(%d)", p.Code)
}
func (p *ClosePacket) Type() PacketType {
	return CLOSE
}
func (p *ClosePacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	if other.Type() != CLOSE {
		return false
	}
	otherClose := other.(*ClosePacket)
	return p.Code == otherClose.Code
}
func (p *ClosePacket) encode() []byte {
	buffer := NewBuffer([]byte{})
	buffer.WriteUInt16(uint16(p.Code))
	return buffer.Bytes()
}

func (p *ClosePacket) decode(data []byte) error {
	buffer := NewBuffer(data)
	code, err := buffer.ReadUInt16()
	if err != nil {
		return err
	}
	p.Code = CloseCode(code)
	return nil
}
