package downloader

import (
	"bufio"
	"os/exec"
	"strings"
)

type Format struct {
	ID   string
	Desc string
}

func ListFormats(url string) ([]Format, error) {
	cmd := exec.Command("yt-dlp", "-F", url)
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
