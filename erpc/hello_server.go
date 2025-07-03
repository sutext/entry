package erpc

import (
	"sutext.github.io/entry/buffer"
)

type HelloReq struct {
	Name string
}

func (req *HelloReq) WriteTo(buf *buffer.Buffer) error {
	buf.WriteString(req.Name)
	return nil
}

func (req *HelloReq) ReadFrom(buf *buffer.Buffer) error {
	name, err := buf.ReadString()
	if err != nil {
		return err
	}
	req.Name = name
	return nil
}

type HelloResp struct {
	Message string
}

func (req *HelloResp) WriteTo(buf *buffer.Buffer) error {
	buf.WriteString(req.Message)
	return nil
}

func (req *HelloResp) ReadFrom(buf *buffer.Buffer) error {
	message, err := buf.ReadString()
	if err != nil {
		return err
	}
	req.Message = message
	return nil
}

type HelloService interface {
	SayHello(req *HelloReq) (*HelloResp, error)
}

type HelloClient struct {
	client *Client
}

func (s *HelloClient) SayHello(req *HelloReq) (*HelloResp, error) {
	return SendMessage[*HelloReq, *HelloResp](s.client, req)
}

type HelloServer struct {
	Service
}

func (s *HelloServer) Name() string {
	return "HelloService"
}
func (s *HelloServer) Handle(req *reqmsg) (res *resmsg, err error) {
	reqObj := &HelloReq{}
	err = reqObj.ReadFrom(buffer.New(req.msg))
	if err != nil {
		return nil, err
	}
	respObj, err := s.SayHello(reqObj)
	if err != nil {
		return nil, err
	}
	buf := buffer.New()
	err = respObj.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	res = &resmsg{seq: req.seq, msg: buf.Bytes()}
	return res, nil
}

func (s *HelloServer) SayHello(req *HelloReq) (*HelloResp, error) {
	resp := &HelloResp{Message: "Hello, " + req.Name + "!"}
	return resp, nil
}
