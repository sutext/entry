package packet

import "sutext.github.io/entry/buffer"

type PingPacket struct{}

func Ping() *PingPacket {
	return &PingPacket{}
}
func (p *PingPacket) Type() PacketType {
	return PING
}
func (p *PingPacket) String() string {
	return PING.String()
}
func (p *PingPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	return PING == other.Type()
}
func (p *PingPacket) WriteTo(w *buffer.Buffer) error {
	return nil
}
func (p *PingPacket) ReadFrom(r *buffer.Buffer) error {
	return nil
}

type PongPacket struct{}

func Pong() *PongPacket {
	return &PongPacket{}
}
func (p *PongPacket) Type() PacketType {
	return PONG
}
func (p *PongPacket) String() string {
	return PONG.String()
}
func (p *PongPacket) Equal(other Packet) bool {
	if other == nil {
		return false
	}
	return PONG == other.Type()
}
func (p *PongPacket) WriteTo(w *buffer.Buffer) error {
	return nil
}
func (p *PongPacket) ReadFrom(r *buffer.Buffer) error {
	return nil
}
