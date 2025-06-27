package code

type Platform uint8

const (
	PlatformUnknown Platform = 0
	PlatformWeb     Platform = 1
	PlatformMini    Platform = 2
	PlatformMobile  Platform = 3
	PlatformDesktop Platform = 4
)

func (p Platform) String() string {
	switch p {
	case PlatformWeb:
		return "web"
	case PlatformMini:
		return "mini"
	case PlatformMobile:
		return "mobile"
	case PlatformDesktop:
		return "desktop"
	default:
		return "unknown"
	}
}
