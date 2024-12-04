package utils

type Win32_UserAccount struct {
	AccountType        uint32
	Description        string
	Disabled           bool
	Domain             string
	FullName           string
	LocalAccount       bool
	Lockout            bool
	Name               string
	PasswordChangeable bool
	PasswordExpires    bool
	PasswordRequired   bool
	SID                string
	SIDType            uint32
	Status             string
}

type KB struct {
	Caption     string
	CSName      string
	Description string
	//FixComments         string
	HotFixID string
	//InstallDate         string
	InstalledBy string
	InstalledOn string
	//Name                string
	//ServicePackInEffect string
}
type Win32_Process struct {
	ProcessID       uint32
	Name            string
	CommandLine     string
	ParentProcessId uint32
	CreationDate    string
	ExecutablePath  string
	Status          string
	ThreadCount     uint32
}
type Win32_NetworkAdapter struct {
	Name         string
	AdapterType  string
	DeviceID     string
	Manufacturer string
	MacAddress   string
	Speed        uint32
	NetEnabled   bool
}
type Win32_BIOS struct {
	Manufacturer      string
	Name              string
	SerialNumber      string
	SMBIOSBIOSVersion string
	Version           string
}

type Win32_Volume struct {
	Capacity    uint64
	DriveType   uint32
	DriveLetter string
	FileSystem  string
	FreeSpace   uint64
	Label       string
	Name        string
}
type Win32_LogicalDisk struct {
	FreeSpace          uint64
	Caption            string
	Size               uint64
	Name               string
	Description        string
	DeviceID           string
	FileSystem         string
	VolumeName         string
	SystemName         string
	ProviderName       string
	VolumeSerialNumber string
}
type Win32_CPU struct {
	Caption                   string
	DeviceID                  string
	Manufacturer              string
	MaxClockSpeed             uint32
	SocketDesignation         string
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	ProcessorId               string
}

type Win32_PhysicalMemory struct {
	Capacity      uint64
	DeviceLocator string
	MemoryType    uint16
	Name          string
	TotalWidth    uint32
	Manufacturer  string
	SerialNumber  string
	PartNumber    string
	Speed         uint32
	Tag           string
}
