package packet

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type DataQos uint8

const (
	DataQosZero DataQos = 0
	DataQosOne  DataQos = 1
)

func (t DataQos) String() string {
	switch t {
	case DataQosZero:
		return "qos0"
	case DataQosOne:
		return "qos1"
	default:
		return "unknown"
	}
}

type DataType byte

const (
	DataTypeText     DataType = 0
	DataTypeJSON     DataType = 1
	DataTypeMsgPack  DataType = 2
	DataTypeProtoBuf DataType = 3
)

func (t DataType) String() string {
	switch t {
	case DataTypeText:
		return "text"
	case DataTypeJSON:
		return "json"
	case DataTypeMsgPack:
		return "msgpack"
	case DataTypeProtoBuf:
		return "protobuf"
	default:
		return "unknown"
	}
}

type DataPacket struct {
	Qos      DataQos
	DataType DataType
	PacketId int64
	Payload  []byte
}

func Data(dataType DataType, packetId int64, payload []byte) *DataPacket {
	return &DataPacket{
		Qos:      DataQosOne,
		DataType: dataType,
		PacketId: packetId,
		Payload:  payload,
	}
}
func Data0(dataType DataType, payload []byte) *DataPacket {
	return &DataPacket{
		Qos:      DataQosZero,
		DataType: dataType,
		PacketId: 0,
		Payload:  payload,
	}
}

func (p *DataPacket) String() string {
	switch p.DataType {
	case DataTypeText:
		return fmt.Sprintf("Text(%s,packetId=%d)", string(p.Payload), p.PacketId)
	case DataTypeJSON:
		var j any
		err := json.Unmarshal(p.Payload, &j)
		if err != nil {
			return fmt.Sprintf("JSONError(%s)", err.Error())
		}
		return fmt.Sprintf("JSON(%s,packetId=%d)", j, p.PacketId)
	case DataTypeMsgPack:
		return fmt.Sprintf("MsgPack(%d bytes,packetId=%d)", len(p.Payload), p.PacketId)
	case DataTypeProtoBuf:
		return fmt.Sprintf("ProtoBuf(%d bytes,packetId=%d)", len(p.Payload), p.PacketId)
	default:
		return fmt.Sprintf("DATA(qos=%s,dataType=%s,packetId=%d, payload=%s)", p.Qos, p.DataType, p.PacketId, string(p.Payload))
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
	return p.PacketId == otherData.PacketId && reflect.DeepEqual(p.Payload, otherData.Payload)
}
func (p *DataPacket) encode() []byte {
	buffer := NewBuffer([]byte{})
	buffer.WriteUInt8(uint8(p.Qos))
	buffer.WriteUInt8(uint8(p.DataType))
	if p.Qos > 0 {
		buffer.WriteInt64(p.PacketId)
	}
	buffer.WriteBytes(p.Payload)
	return buffer.Bytes()
}

func (p *DataPacket) decode(data []byte) error {
	buffer := NewBuffer(data)
	qos, err := buffer.ReadUInt8()
	if err != nil {
		return err
	}
	dataType, err := buffer.ReadUInt8()
	if err != nil {
		return err
	}
	p.Qos = DataQos(qos)
	p.DataType = DataType(dataType)
	if qos > 0 {
		packetId, err := buffer.ReadInt64()
		if err != nil {
			return err
		}
		p.PacketId = packetId
	}
	payload, err := buffer.Read()
	if err != nil {
		return err
	}
	p.Payload = payload
	return nil
}

type DataAckPacket struct {
	PacketId int64
}

func DataAck(packetId int64) Packet {
	return &DataAckPacket{PacketId: packetId}
}
func (p *DataAckPacket) String() string {
	return fmt.Sprintf("DATAACK(packetId=%d)", p.PacketId)
}
func (p *DataAckPacket) Type() PacketType {
	return DATAACK
}
func (p *DataAckPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	if p.Type() != other.Type() {
		return false
	}
	otherAck := other.(*DataAckPacket)
	return p.PacketId == otherAck.PacketId
}
func (p *DataAckPacket) encode() []byte {
	buffer := NewBuffer([]byte{})
	buffer.WriteInt64(p.PacketId)
	return buffer.Bytes()
}

func (p *DataAckPacket) decode(data []byte) error {
	buffer := NewBuffer(data)
	packetId, err := buffer.ReadInt64()
	if err != nil {
		return err
	}
	p.PacketId = packetId
	return nil
}
