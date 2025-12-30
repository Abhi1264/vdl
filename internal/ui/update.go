package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Abhi1264/vidforge/internal/config"
	"github.com/Abhi1264/vidforge/internal/downloader"
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		specialKeys := map[string]bool{
			"ctrl+c": true,
			"q":      true,
			"p":      true,
			"up":     true,
			"k":      true,
			"down":   true,
			"j":      true,
			"?":      true,
			"enter":  true,
		}

		if m.showDownloadPath {
			specialKeys["d"] = true
		} else if !m.showHelp && !m.showProfiles {
			specialKeys["s"] = true
			specialKeys["f"] = true
			specialKeys["d"] = true
			specialKeys["1"] = true
			specialKeys["2"] = true
			specialKeys["3"] = true
			specialKeys["4"] = true
			specialKeys["5"] = true
			specialKeys["6"] = true
		} else if m.showProfiles {
			specialKeys["f"] = true
			specialKeys["1"] = true
			specialKeys["2"] = true
			specialKeys["3"] = true
			specialKeys["4"] = true
			specialKeys["5"] = true
			specialKeys["6"] = true
		}

		if !specialKeys[key] {
			if m.showDownloadPath {
				m.downloadPathInput, cmd = m.downloadPathInput.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.input, cmd = m.input.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "p":
			if !m.showDownloadPath {
				if job := m.selectedJob(); job != nil {
					m.manager.Cancel(m.selectedID)
					job.status = "cancelled"
					job.done = true
				}
			}
		case "up", "k":
			if m.showDownloadPath {
				return m, tea.Batch(cmds...)
			}
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
			if m.showDownloadPath {
				return m, tea.Batch(cmds...)
			}
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
				m.showDownloadPath = false
			}
		case "s":
			if !m.showDownloadPath {
				m.sponsorBlock = !m.sponsorBlock
			}
		case "f":
			if !m.showDownloadPath {
				m.showProfiles = !m.showProfiles
				if m.showProfiles {
					m.showHelp = false
				}
			}
		case "d":
			if m.showDownloadPath {
				m.showDownloadPath = false
				m.downloadPathInput.Blur()
				m.input.Focus()
			} else if !m.showHelp && !m.showProfiles {
				m.showDownloadPath = true
				m.input.Blur()
				m.downloadPathInput.Focus()
				m.downloadPathInput.SetValue(m.downloadPath)
			}
			return m, tea.Batch(cmds...)
		case "1", "2", "3", "4", "5", "6":
			if m.showDownloadPath {
				return m, tea.Batch(cmds...)
			}
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
				return m, tea.Batch(cmds...)
			}
			if m.showDownloadPath {
				newPath := m.downloadPathInput.Value()
				if newPath != "" {
					// Expand ~ if present
					if len(newPath) > 0 && newPath[0] == '~' {
						homeDir, err := os.UserHomeDir()
						if err == nil {
							newPath = filepath.Join(homeDir, newPath[1:])
						}
					}
					m.downloadPath = newPath
					cfg := config.GetConfig()
					cfg.SetDownloadPath(newPath)
				}
				m.showDownloadPath = false
				m.downloadPathInput.Blur()
				m.input.Focus()
				return m, tea.Batch(cmds...)
			}
			if m.showProfiles {
				profiles := downloader.GetProfiles()
				if m.profileIndex >= 0 && m.profileIndex < len(profiles) {
					m.profile = &profiles[m.profileIndex]
					m.showProfiles = false
				}
				return m, tea.Batch(cmds...)
			}
			url := m.input.Value()
			if url == "" {
				return m, tea.Batch(cmds...)
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
				OutputPath:   m.downloadPath,
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

	cmds = append(cmds, listen(m.progressCh))
	return m, tea.Batch(cmds...)
}
