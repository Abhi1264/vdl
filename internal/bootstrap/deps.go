package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Ensure(cmd string) error {
	if _, err := exec.LookPath(cmd); err == nil {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return installMacOS(cmd)
	case "linux":
		return installLinux(cmd)
	case "windows":
		return installWindows(cmd)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func installMacOS(cmd string) error {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return fmt.Errorf("Homebrew not found. Please install Homebrew first: https://brew.sh")
	}

	var pkg string
	switch cmd {
	case "yt-dlp":
		pkg = "yt-dlp"
	case "ffmpeg":
		pkg = "ffmpeg"
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	installCmd := exec.Command(brewPath, "install", pkg)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	return installCmd.Run()
}

func installLinux(cmd string) error {
	var installCmd *exec.Cmd

	if _, err := exec.LookPath("apt-get"); err == nil {
		switch cmd {
		case "yt-dlp":
			installCmd = exec.Command("sudo", "apt-get", "update")
			if err := installCmd.Run(); err != nil {
				return err
			}
			installCmd = exec.Command("sudo", "apt-get", "install", "-y", "yt-dlp")
		case "ffmpeg":
			installCmd = exec.Command("sudo", "apt-get", "install", "-y", "ffmpeg")
		}
	} else if _, err := exec.LookPath("dnf"); err == nil {
		switch cmd {
		case "yt-dlp":
			installCmd = exec.Command("sudo", "dnf", "install", "-y", "yt-dlp")
		case "ffmpeg":
			installCmd = exec.Command("sudo", "dnf", "install", "-y", "ffmpeg")
		}
	} else if _, err := exec.LookPath("pacman"); err == nil {
		switch cmd {
		case "yt-dlp":
			installCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "yt-dlp")
		case "ffmpeg":
			installCmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "ffmpeg")
		}
	} else {
		return fmt.Errorf("no supported package manager found (apt-get, dnf, or pacman)")
	}

	if installCmd == nil {
		return fmt.Errorf("unknown command: %s", cmd)
	}

	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	return installCmd.Run()
}

func installWindows(cmd string) error {
	exeDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable directory: %w", err)
	}
	appDir := filepath.Dir(exeDir)

	var url string
	var filename string

	switch cmd {
	case "yt-dlp":
		url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
		filename = "yt-dlp.exe"
	case "ffmpeg":
		arch := "win64"
		if runtime.GOARCH == "386" {
			arch = "win32"
		}
		url = fmt.Sprintf("https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-%s.zip", arch)
		filename = "ffmpeg.zip"
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	if cmd == "yt-dlp" {
		destPath := filepath.Join(appDir, filename)
		if err := downloadFile(url, destPath); err != nil {
			return fmt.Errorf("failed to download %s: %w", cmd, err)
		}
		return nil
	}

	return fmt.Errorf("Windows installation for %s not yet fully implemented", cmd)
}

func downloadFile(url, dest string) error {
	if runtime.GOOS == "windows" {
		psScript := fmt.Sprintf(`
			$ProgressPreference = 'SilentlyContinue'
			Invoke-WebRequest -Uri "%s" -OutFile "%s"
		`, url, dest)
		cmd := exec.Command("powershell", "-Command", psScript)
		return cmd.Run()
	}

	curlPath, err := exec.LookPath("curl")
	if err != nil {
		wgetPath, err := exec.LookPath("wget")
		if err != nil {
			return fmt.Errorf("neither curl nor wget found")
		}
		cmd := exec.Command(wgetPath, "-O", dest, url)
		return cmd.Run()
	}

	cmd := exec.Command(curlPath, "-L", "-o", dest, url)
	return cmd.Run()
}
