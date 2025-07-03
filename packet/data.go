package packet

import (
	"fmt"
	"reflect"

	"sutext.github.io/entry/buffer"
)

type DataEncoding byte

const (
	DataText   DataEncoding = 0
	DataJSON   DataEncoding = 1
	DataBinary DataEncoding = 2
)

func (t DataEncoding) String() string {
	switch t {
	case DataText:
		return "Text"
	case DataJSON:
		return "JSON"
	case DataBinary:
		return "Binary"
	default:
		return "Custom"
	}
}

type DataPacket struct {
	Encoding DataEncoding
	Payload  []byte
}

func Data(dataType DataEncoding, payload []byte) *DataPacket {
	return &DataPacket{
		Encoding: dataType,
		Payload:  payload,
	}
}

func (p *DataPacket) String() string {
	switch p.Encoding {
	case DataText:
		return fmt.Sprintf("DATA{%s, Payload: %s}", p.Encoding, string(p.Payload))
	case DataJSON:
		return fmt.Sprintf("DATA{%s, Payload: %s}", p.Encoding, string(p.Payload))
	case DataBinary:
		return fmt.Sprintf("DATA{%s, Payload: %d bytes}", p.Encoding, len(p.Payload))
	default:
		return fmt.Sprintf("DATA{%s, Payload: %d bytes}", p.Encoding, len(p.Payload))
	}
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
	return p.Encoding == otherData.Encoding && reflect.DeepEqual(p.Payload, otherData.Payload)
}

func (p *DataPacket) WriteTo(buf *buffer.Buffer) error {
	buf.WriteUInt8(uint8(p.Encoding))
	buf.WriteBytes(p.Payload)
	return nil
}

func (p *DataPacket) ReadFrom(buf *buffer.Buffer) error {
	dataType, err := buf.ReadUInt8()
	if err != nil {
		return err
	}
	p.Encoding = DataEncoding(dataType)
	p.Payload, err = buf.ReadAll()
	return err
}
