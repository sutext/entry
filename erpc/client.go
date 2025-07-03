package erpc

import (
	"sync/atomic"

	"sutext.github.io/entry/buffer"
	"sutext.github.io/entry/client"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
)

type MessageChannel struct {
	resp Message
	done chan struct{}
}
type Client struct {
	cli  *client.Client
	seq  atomic.Int64
	reqs *safe.Map[int64, MessageChannel]
}

func Dail(host, prot string) (*Client, error) {
	cfg := client.NewConfig()
	cfg.Host = host
	cfg.Port = prot
	cli := client.New(cfg)
	cli.Connect(nil)
	c := &Client{
		cli:  cli,
		reqs: safe.NewMap[map[int64]MessageChannel](),
	}
	cli.HandleData(c.handleData)
	return c, nil
}

func SendMessage[Req Message, Resp Message](c *Client, req Req) (resp Resp, err error) {
	buf := buffer.New()
	seq := c.seq.Add(1)
	reqmsg := reqmsg{
		seq: seq,
		msg: buf.Bytes(),
	}
	buf.WriteUInt8(0)
	reqmsg.WriteTo(buf)
	err = c.cli.SendData(buf.Bytes())
	if err != nil {
		return resp, err
	}
	ch := MessageChannel{resp: resp, done: make(chan struct{})}
	c.reqs.Set(seq, ch)
	<-ch.done
	resp = ch.resp.(Resp)
	return resp, nil
}
func (c *Client) handleData(p *packet.DataPacket) error {
	buf := buffer.New(p.Payload)
	flag, err := buf.ReadUInt8()
	if err != nil {
		return err
	}
	if flag > 0 {
		res := &resmsg{}
		err = res.ReadFrom(buf)
		if err != nil {
			return err
		}
		return c.handleResponse(res)
	} else {
		req := &reqmsg{}
		err = req.ReadFrom(buf)
		if err != nil {
			return err
		}
		return c.handleRequest(req)
	}

}
func (c *Client) handleRequest(req *reqmsg) error {
	return nil
}
func (c *Client) handleResponse(res *resmsg) error {
	ch, ok := c.reqs.Get(res.seq)
	if ok {
		ch.resp.ReadFrom(buffer.New(res.msg))
		close(ch.done)
	}
	return nil
}
