package downloader

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Abhi1264/vidforge/internal/bootstrap"
)

type Format struct {
	ID   string
	Desc string
}

func ListFormats(url string) ([]Format, error) {
	ytdlpPath, err := bootstrap.GetCommandPath("yt-dlp")
	if err != nil {
		return nil, fmt.Errorf("yt-dlp not found: %w", err)
	}

	cmd := exec.Command(ytdlpPath, "-F", url)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var formats []Format
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] >= '0' && line[0] <= '9' {
			parts := strings.Fields(line)
			formats = append(formats, Format{
				ID:   parts[0],
				Desc: strings.Join(parts[1:], " "),
			})
		}
	}
	return formats, cmd.Wait()
}
