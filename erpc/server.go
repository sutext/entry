package erpc

import (
	"fmt"

	"sutext.github.io/entry/buffer"
	"sutext.github.io/entry/packet"
	"sutext.github.io/entry/safe"
	"sutext.github.io/entry/server"
)

type Service interface {
	Name() string
	Handle(req *reqmsg) (res *resmsg, err error)
}

type Server struct {
	impl     *server.Server
	services *safe.Map[string, Service]
}

func NewServer() *Server {
	conf := server.NewConfig()
	impl := server.New(conf)
	s := &Server{
		impl:     impl,
		services: safe.NewMap(map[string]Service{}),
	}
	s.impl.HandleData(s.handleData)
	return s
}
func (s *Server) Register(svc Service) {
	s.services.Set(svc.Name(), svc)
}

func (c *Server) handleData(p *packet.DataPacket) (*packet.DataPacket, error) {
	buf := buffer.New(p.Payload)
	flag, err := buf.ReadUInt8()
	if err != nil {
		return nil, err
	}
	if flag > 0 {
		return nil, fmt.Errorf("unsupported flag %d", flag)
	} else {
		req := &reqmsg{}
		err = req.ReadFrom(buf)
		if err != nil {
			return nil, err
		}
		service, ok := c.services.Get(req.service)
		if !ok {
			return nil, fmt.Errorf("service %s not found", req.service)
		}
		res, err := service.Handle(req)
		if err != nil {
			return nil, err
		}
		buf := buffer.New()
		buf.WriteUInt8(1)
		res.WriteTo(buf)
		return packet.Data(packet.DataBinary, buf.Bytes()), nil
	}

}
