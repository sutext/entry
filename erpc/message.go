package erpc

import "sutext.github.io/entry/buffer"

type Message interface {
	buffer.WriteTo
	buffer.ReadFrom
}

type reqmsg struct {
	seq     int64
	service string
	method  string
	msg     []byte
}

func (m *reqmsg) WriteTo(buf *buffer.Buffer) error {
	buf.WriteVarint(m.seq)
	buf.WriteString(m.service)
	buf.WriteString(m.method)
	buf.WriteBytes(m.msg)
	return nil
}
func (m *reqmsg) ReadFrom(buf *buffer.Buffer) error {
	seq, err := buf.ReadVarint()
	if err != nil {
		return err
	}
	m.seq = seq
	m.service, err = buf.ReadString()
	if err != nil {
		return err
	}
	m.method, err = buf.ReadString()
	if err != nil {
		return err
	}
	m.msg, err = buf.ReadAll()
	return err
}

type resmsg struct {
	seq int64
	msg []byte
}

func (m *resmsg) WriteTo(buf *buffer.Buffer) error {
	buf.WriteVarint(m.seq)
	buf.WriteBytes(m.msg)
	return nil
}
func (m *resmsg) ReadFrom(buf *buffer.Buffer) error {
	seq, err := buf.ReadVarint()
	if err != nil {
		return err
	}
	m.seq = seq
	m.msg, err = buf.ReadAll()
	return err
}
