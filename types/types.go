package types

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

type Platform uint8

const (
	PlatformUnknown Platform = 0
	PlatformWeb     Platform = 1
	PlatformMini    Platform = 2
	PlatformMobile  Platform = 3
	PlatformDesktop Platform = 4
)
