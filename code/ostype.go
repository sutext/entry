package code

type OSType uint8

const (
	OSTypeUnknown OSType = 0
	OSTypeIOS     OSType = 1
	OSTypeAndroid OSType = 2
	OSTypeHarmony OSType = 3
	OSTypeWindows OSType = 4
	OSTypeMacOS   OSType = 5
	OSTypeLinux   OSType = 6
)

func (t OSType) String() string {
	switch t {
	case OSTypeUnknown:
		return "Unknown"
	case OSTypeIOS:
		return "iOS"
	case OSTypeAndroid:
		return "Android"
	case OSTypeHarmony:
		return "Harmony"
	case OSTypeWindows:
		return "Windows"
	case OSTypeMacOS:
		return "MacOS"
	case OSTypeLinux:
		return "Linux"
	default:
		return "Unknown"
	}
}

func (t OSType) Error() string {
	return t.String()
}
