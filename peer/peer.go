package peer

type Server interface {
	SendData(data []byte) error
	KickOut(cid string) error
}

func Register(cid string, ip string) error {
	return nil
}
func Unregister(cid string) error {
	return nil
}
