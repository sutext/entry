package client

type Status uint8

const (
	StatusUnknown Status = iota
	StatusClosed
	StatusOpened
	StatusOpening
	StatusClosing
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "Unknown"
	case StatusClosed:
		return "Closed"
	case StatusOpened:
		return "Opened"
	case StatusOpening:
		return "Opening"
	case StatusClosing:
		return "Closing"
	default:
		return "Unknown"
	}
}
