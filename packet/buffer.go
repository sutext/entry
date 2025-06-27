package packet

import (
	"encoding/binary"
)

type buffer struct {
	pos int
	buf []byte
}

func newBuffer(buf []byte) *buffer {
	return &buffer{pos: 0, buf: buf}
}
func (b *buffer) bytes() []byte {
	return b.buf
}
func (b *buffer) writeUInt8(i uint8) {
	b.buf = append(b.buf, i)
}
func (b *buffer) writeBytes(p []byte) {
	b.buf = append(b.buf, p...)
}
func (b *buffer) writeUInt16(i uint16) {
	b.buf = binary.BigEndian.AppendUint16(b.buf, i)
}
func (b *buffer) writeUInt32(i uint32) {
	b.buf = binary.BigEndian.AppendUint32(b.buf, i)
}
func (b *buffer) writeInt64(i int64) {
	b.writeUInt64(uint64(i))
}
func (b *buffer) writeUInt64(i uint64) {
	b.buf = binary.BigEndian.AppendUint64(b.buf, i)
}
func (b *buffer) writeVarint(i int64) {
	b.buf = binary.AppendVarint(b.buf, i)
}
func (b *buffer) writeString(s string) {
	bytes := []byte(s)
	b.writeVarint(int64(len(bytes)))
	b.writeBytes(bytes)
}

func (b *buffer) readAll() ([]byte, error) {
	if b.pos >= len(b.buf) {
		return nil, ErrBufferTooShort
	}
	p := b.buf[b.pos:]
	b.pos = len(b.buf)
	return p, nil
}
func (b *buffer) readUInt8() (uint8, error) {
	if b.pos+1 >= len(b.buf) {
		return 0, ErrBufferTooShort
	}
	i := b.buf[b.pos]
	b.pos++
	return i, nil
}
func (b *buffer) readBytes(l int) ([]byte, error) {
	if b.pos+l > len(b.buf) {
		return nil, ErrBufferTooShort
	}
	p := b.buf[b.pos : b.pos+l]
	b.pos += l
	return p, nil
}

func (b *buffer) readUInt16() (uint16, error) {
	bytes, err := b.readBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bytes), nil
}

func (b *buffer) readUInt32() (uint32, error) {
	bytes, err := b.readBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(bytes), nil
}

func (b *buffer) readUInt64() (uint64, error) {
	bytes, err := b.readBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(bytes), nil
}
func (b *buffer) readInt64() (int64, error) {
	u, err := b.readUInt64()
	if err != nil {
		return 0, err
	}
	return int64(u), nil
}

func (b *buffer) readVarint() (int64, error) {
	varint, len := binary.Varint(b.buf[b.pos:])
	if len < 0 {
		return 0, ErrVarintOverflow
	}
	if len == 0 {
		return 0, ErrBufferTooShort
	}
	b.pos += len
	return varint, nil
}
func (b *buffer) readString() (string, error) {
	l, err := b.readVarint()
	if err != nil {
		return "", err
	}
	bytes, err := b.readBytes(int(l))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
