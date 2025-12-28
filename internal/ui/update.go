package ui

import (
	"fmt"
	"github.com/Abhi1264/vdl/internal/downloader"
	tea "github.com/charmbracelet/bubbletea"
)

type progressMsg downloader.Progress

func listen(ch <-chan downloader.Progress) tea.Cmd {
	return func() tea.Msg {
		if p, ok := <-ch; ok {
			return progressMsg(p)
		}
		return nil
	}
}

func (m Model) Init() tea.Cmd {
	return listen(m.progressCh)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if job := m.selectedJob(); job != nil {
				m.manager.Cancel(m.selectedID)
				job.status = "cancelled"
				job.done = true
			}
		case "up", "k":
			if m.showProfiles {
				profiles := downloader.GetProfiles()
				if len(profiles) > 0 {
					m.profileIndex--
					if m.profileIndex < 0 {
						m.profileIndex = len(profiles) - 1
					}
					m.profile = &profiles[m.profileIndex]
				}
				return m, cmd
			}
			if len(m.jobs) > 0 {
				ids := m.getSortedIDs()
				if len(ids) > 0 {
					currentIdx := m.findIDIndex(ids, m.selectedID)
					if currentIdx > 0 {
						m.selectedID = ids[currentIdx-1]
					} else {
						m.selectedID = ids[len(ids)-1]
					}
				}
			}
		case "down", "j":
			if m.showProfiles {
				profiles := downloader.GetProfiles()
				if len(profiles) > 0 {
					m.profileIndex++
					if m.profileIndex >= len(profiles) {
						m.profileIndex = 0
					}
					m.profile = &profiles[m.profileIndex]
				}
				return m, cmd
			}
			if len(m.jobs) > 0 {
				ids := m.getSortedIDs()
				if len(ids) > 0 {
					currentIdx := m.findIDIndex(ids, m.selectedID)
					if currentIdx < len(ids)-1 {
						m.selectedID = ids[currentIdx+1]
					} else {
						m.selectedID = ids[0]
					}
				}
			}
		case "?":
			m.showHelp = !m.showHelp
			if m.showHelp {
				m.showProfiles = false
			}
		case "s":
			m.sponsorBlock = !m.sponsorBlock
		case "f":
			m.showProfiles = !m.showProfiles
			if m.showProfiles {
				m.showHelp = false
			}
		case "1", "2", "3", "4", "5", "6":
			if m.showProfiles {
				idx := int(msg.String()[0] - '1')
				profiles := downloader.GetProfiles()
				if idx >= 0 && idx < len(profiles) {
					m.profileIndex = idx
					m.profile = &profiles[idx]
					m.showProfiles = false
				}
				return m, cmd
			}

		case "enter":
			if m.showHelp {
				return m, cmd
			}
			if m.showProfiles {
				profiles := downloader.GetProfiles()
				if m.profileIndex >= 0 && m.profileIndex < len(profiles) {
					m.profile = &profiles[m.profileIndex]
					m.showProfiles = false
				}
				return m, cmd
			}
			url := m.input.Value()
			if url == "" {
				return m, cmd
			}

			id := m.nextID
			m.nextID++

			isYouTube := downloader.IsYouTubeURL(url)
			sponsorBlock := m.sponsorBlock && isYouTube

			m.jobs[id] = &jobState{
				url:    url,
				status: "downloading",
			}

			if len(m.jobs) == 1 {
				m.selectedID = id
			}

			m.manager.Submit(downloader.Job{
				ID:           id,
				URL:          url,
				Profile:      m.profile,
				SponsorBlock: sponsorBlock,
				Resume:       true,
			})

			m.input.Reset()
		}

	case progressMsg:
		job := m.jobs[msg.ID]
		if job == nil {
			return m, listen(m.progressCh)
		}

		if msg.Err != nil {
			job.status = fmt.Sprintf("error: %v", msg.Err)
			job.done = true
		} else if msg.Done {
			job.progress = 1.0
			job.status = "done"
			job.done = true
		} else {
			job.progress = msg.Percent
		}
		return m, listen(m.progressCh)
	}

	return m, tea.Batch(cmd, listen(m.progressCh))
}
