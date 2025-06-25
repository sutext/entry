package packet

import "errors"

var (
	ErrNegativeLength     = errors.New("negative length")
	ErrVarintOverflow     = errors.New("varint overflow")
	ErrBufferTooShort     = errors.New("buffer too short")
	ErrUnkownPacketType   = errors.New("unknown packet type")
	ErrPacketSizeTooLarge = errors.New("packet size too large")
)
