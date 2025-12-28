package ui

import (
	"github.com/Abhi1264/vdl/internal/downloader"
)

type formatItem downloader.Format

func (f formatItem) Title() string       { return f.ID }
func (f formatItem) Description() string { return f.Desc }
func (f formatItem) FilterValue() string { return f.ID }
