package ui

import (
	"github.com/Abhi1264/vdl/internal/downloader"
	"github.com/charmbracelet/bubbles/textinput"
)

type jobState struct {
	url      string
	progress float64
	status   string
	done     bool
}

type Model struct {
	input        textinput.Model
	jobs         map[int]*jobState
	nextID       int
	selectedID   int
	manager      *downloader.Manager
	progressCh   <-chan downloader.Progress
	showHelp     bool
	showProfiles bool
	profileIndex int
	profile      *downloader.Profile
	sponsorBlock bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "https://..."
	ti.Focus()
	ti.CharLimit = 2048
	ti.Width = 60

	mgr := downloader.NewManager(3)
	profiles := downloader.GetProfiles()
	defaultProfile := downloader.GetDefaultProfile()
	profileIndex := 0
	for i, p := range profiles {
		if p.Quality == defaultProfile.Quality {
			profileIndex = i
			break
		}
	}

	return Model{
		input:        ti,
		jobs:         make(map[int]*jobState),
		manager:      mgr,
		progressCh:   mgr.Updates(),
		profile:      defaultProfile,
		profileIndex: profileIndex,
		sponsorBlock: true,
	}
}

func (m Model) selectedJob() *jobState {
	return m.jobs[m.selectedID]
}

func (m Model) getSortedIDs() []int {
	ids := make([]int, 0, len(m.jobs))
	for id := range m.jobs {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids)-1; i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] > ids[j] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}
	return ids
}

func (m Model) findIDIndex(ids []int, target int) int {
	for i, id := range ids {
		if id == target {
			return i
		}
	}
	if len(ids) > 0 {
		return 0
	}
	return -1
}
