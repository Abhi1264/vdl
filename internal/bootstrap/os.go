package bootstrap

import "runtime"

type OSInfo struct {
	OS   string
	Arch string
}

func DetectOS() OSInfo {
	return OSInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsIntel() bool {
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "386"
}

func IsARM() bool {
	return runtime.GOARCH == "arm64" || runtime.GOARCH == "arm"
}
