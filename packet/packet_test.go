package packet

import (
	"fmt"
	"io"
	"testing"
)

func TestPacket(t *testing.T) {
	identity := Identity{}
	testp(t, Connect(&identity))
	testp(t, Connack(0))
	testp(t, SmallData())
	testp(t, BigData())
	testp(t, Ping())
	testp(t, Pong())
	testp(t, Close(0))
}
func SmallData() Packet {
	return Data([]byte("hello world"))
}
func BigData() Packet {
	return Data([]byte("hellfasdfhellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;ohellfasdfsafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;osafsadfasdfasdfasdfsdafasdfasdfsd;afjjjjdslfjasl;dkfjasdlfjasdl;kfjasl;dfjsadlk;fjasl;dkfjasldk;fjlskd;ajflasdkfjlaksdjfalsd;fjalsd;kfjsadl;kfjasl;dkfjasdl;kfjalsdkfjals;dkfjasld;kfjals;dkfjsal;dkfjalsd;kfjlas;dkfjalsdkfjsld;fjals;dfjals;fjaldsk;o"))
}
func testp(t *testing.T, p Packet) {
	rw := &ReadWriter{}
	err := WritePacket(rw, p)
	if err != nil {
		t.Error(err)
	}
	newp, err := ReadPacket(rw)
	if err != nil {
		t.Error(err)
	}
	if !p.Equal(newp) {
		fmt.Printf("old packet: %v\n", p)
		fmt.Printf("new packet: %v\n", newp)
		t.Error("data packet not equal")
	}
}

type ReadWriter struct {
	data []byte
}

func (w *ReadWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *ReadWriter) Read(p []byte) (n int, err error) {
	l := len(p)
	if l == 0 {
		return 0, nil
	}
	if len(w.data) == 0 {
		return 0, io.EOF
	}
	if l < len(w.data) {
		n = copy(p, w.data[:l])
		w.data = w.data[n:]
		return n, nil
	} else {
		n = copy(p, w.data)
		w.data = nil
		return n, nil
	}
}
