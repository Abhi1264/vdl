package downloader

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/Abhi1264/vidforge/internal/bootstrap"
)

var percentRe = regexp.MustCompile(`(\d{1,3}\.\d)%`)

type Progress struct {
	ID      int
	Percent float64
	Done    bool
	Err     error
}

type Job struct {
	ID           int
	URL          string
	Profile      *Profile
	SponsorBlock bool
	Resume       bool
	OutputPath   string
}

func (j Job) Run(ctx context.Context, out chan<- Progress) {
	args := []string{
		"--newline",
		"--progress-template",
		"%(progress._percent_str)s",
	}

	if j.Resume {
		args = append(args, "--continue", "--no-part")
	}

	if j.Profile != nil {
		args = append(args, j.Profile.Flags...)
	}

	if j.SponsorBlock {
		args = append(args,
			"--sponsorblock-remove", "sponsor,selfpromo,interaction,intro,outro",
		)
	}

	if j.OutputPath != "" {
		if err := os.MkdirAll(j.OutputPath, 0755); err != nil {
			out <- Progress{ID: j.ID, Err: err}
			return
		}
		outputTemplate := "%(title)s.%(ext)s"
		outputPath := filepath.Join(j.OutputPath, outputTemplate)
		args = append(args, "-o", outputPath)
	}

	args = append(args, j.URL)

	ytdlpPath, err := bootstrap.GetCommandPath("yt-dlp")
	if err != nil {
		out <- Progress{ID: j.ID, Err: fmt.Errorf("yt-dlp not found: %w", err)}
		return
	}

	cmd := exec.CommandContext(ctx, ytdlpPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		out <- Progress{ID: j.ID, Err: err}
		return
	}

	if err := cmd.Start(); err != nil {
		out <- Progress{ID: j.ID, Err: err}
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			cmd.Process.Kill()
			out <- Progress{ID: j.ID, Done: true, Err: ctx.Err()}
			return
		default:
		}
		line := scanner.Text()
		match := percentRe.FindStringSubmatch(line)
		if len(match) == 2 {
			val, _ := strconv.ParseFloat(match[1], 64)
			out <- Progress{ID: j.ID, Percent: val / 100}
		}
	}

	err = cmd.Wait()
	if ctx.Err() != nil {
		out <- Progress{ID: j.ID, Done: true, Err: ctx.Err()}
	} else {
		out <- Progress{ID: j.ID, Done: true, Err: err}
	}
}
