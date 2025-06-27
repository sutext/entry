package code

type CloseCode uint16

const (
	CloseNormal CloseCode = 1000
	CloseGoaway CloseCode = 1001
)

type ConnectCode uint16

const (
	ConnectAccepted ConnectCode = 0
)
