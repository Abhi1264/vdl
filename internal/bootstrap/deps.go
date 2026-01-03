package bootstrap

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	ytdlpDarwinAMD64  []byte
	ytdlpDarwinARM64  []byte
	ytdlpLinuxAMD64   []byte
	ytdlpLinuxARM64   []byte
	ytdlpWindowsAMD64 []byte
	ytdlpWindowsARM64 []byte

	ffmpegDarwinAMD64  []byte
	ffmpegDarwinARM64  []byte
	ffmpegLinuxAMD64   []byte
	ffmpegLinuxARM64   []byte
	ffmpegWindowsAMD64 []byte
	ffmpegWindowsARM64 []byte
)

var extractedBinaries = make(map[string]string)

func getEmbeddedBinary(cmd string) ([]byte, error) {
	platform := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)

	var binaryData []byte
	var found bool

	switch cmd {
	case "yt-dlp":
		switch platform {
		case "darwin_amd64":
			binaryData, found = ytdlpDarwinAMD64, len(ytdlpDarwinAMD64) > 0
		case "darwin_arm64":
			binaryData, found = ytdlpDarwinARM64, len(ytdlpDarwinARM64) > 0
		case "linux_amd64":
			binaryData, found = ytdlpLinuxAMD64, len(ytdlpLinuxAMD64) > 0
		case "linux_arm64":
			binaryData, found = ytdlpLinuxARM64, len(ytdlpLinuxARM64) > 0
		case "windows_amd64":
			binaryData, found = ytdlpWindowsAMD64, len(ytdlpWindowsAMD64) > 0
		case "windows_arm64":
			binaryData, found = ytdlpWindowsARM64, len(ytdlpWindowsARM64) > 0
		}
	case "ffmpeg":
		switch platform {
		case "darwin_amd64":
			binaryData, found = ffmpegDarwinAMD64, len(ffmpegDarwinAMD64) > 0
		case "darwin_arm64":
			binaryData, found = ffmpegDarwinARM64, len(ffmpegDarwinARM64) > 0
		case "linux_amd64":
			binaryData, found = ffmpegLinuxAMD64, len(ffmpegLinuxAMD64) > 0
		case "linux_arm64":
			binaryData, found = ffmpegLinuxARM64, len(ffmpegLinuxARM64) > 0
		case "windows_amd64":
			binaryData, found = ffmpegWindowsAMD64, len(ffmpegWindowsAMD64) > 0
		case "windows_arm64":
			binaryData, found = ffmpegWindowsARM64, len(ffmpegWindowsARM64) > 0
		}
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}

	if !found {
		return nil, fmt.Errorf("embedded binary for %s not found for platform %s", cmd, platform)
	}

	return binaryData, nil
}

func Ensure(cmd string) error {
	if _, err := exec.LookPath(cmd); err == nil {
		return nil
	}

	if path, ok := extractedBinaries[cmd]; ok {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		delete(extractedBinaries, cmd)
	}

	return extractBinary(cmd)
}

func extractBinary(cmd string) error {
	binaryData, err := getEmbeddedBinary(cmd)
	if err != nil {
		return fmt.Errorf("failed to get embedded binary: %w", err)
	}

	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = cmd + ".exe"
	} else {
		binaryName = cmd
	}

	exeDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable directory: %w", err)
	}
	appDir := filepath.Dir(exeDir)

	binDir := filepath.Join(appDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	destPath := filepath.Join(binDir, binaryName)
	if err := os.WriteFile(destPath, binaryData, 0755); err != nil {
		return fmt.Errorf("failed to write binary: %w", err)
	}

	extractedBinaries[cmd] = destPath
	return nil
}

func GetCommandPath(cmd string) (string, error) {
	if path, ok := extractedBinaries[cmd]; ok {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		delete(extractedBinaries, cmd)
	}

	if path, err := exec.LookPath(cmd); err == nil {
		return path, nil
	}

	if err := extractBinary(cmd); err == nil {
		if path, ok := extractedBinaries[cmd]; ok {
			return path, nil
		}
	}

	return "", fmt.Errorf("command %s not found in PATH or embedded binaries", cmd)
}
