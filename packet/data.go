package packet

import (
	"fmt"
	"reflect"

	"sutext.github.io/entry/buffer"
)

type DataPacket struct {
	Group   string
	Payload []byte
}

func NewData(payload []byte) *DataPacket {
	return &DataPacket{
		Payload: payload,
	}
}

func (p *DataPacket) String() string {
	return fmt.Sprintf("DATA(%d bytes)", len(p.Payload))
}
func (p *DataPacket) Type() PacketType {
	return DATA
}
func (p *DataPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	if p.Type() != other.Type() {
		return false
	}
	otherData := other.(*DataPacket)
	return p.Group == otherData.Group && reflect.DeepEqual(p.Payload, otherData.Payload)
}

func (p *DataPacket) WriteTo(buf *buffer.Buffer) error {
	buf.WriteString(p.Group)
	buf.WriteBytes(p.Payload)
	return nil
}

func (p *DataPacket) ReadFrom(buf *buffer.Buffer) error {
	group, err := buf.ReadString()
	if err != nil {
		return err
	}
	payload, err := buf.ReadAll()
	if err != nil {
		return err
	}
	p.Group = group
	p.Payload = payload
	return nil
}
