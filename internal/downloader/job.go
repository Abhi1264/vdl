package downloader

import (
	"bufio"
	"context"
	"os/exec"
	"regexp"
	"strconv"
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

	args = append(args, j.URL)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

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
