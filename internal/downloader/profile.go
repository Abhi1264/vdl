package downloader

type Profile struct {
	Name        string
	Specs       string
	Quality     string
	Flags       []string
	Description string
}

var profiles = []Profile{
	{
		Name:        "Best Quality",
		Specs:       "4K/8K • VP9/AV1 + Opus • MKV",
		Quality:     "best",
		Description: "Best Quality — 4K/8K • VP9/AV1 + Opus • MKV (best)",
		Flags: []string{
			"-f", "bestvideo+bestaudio/best",
			"--merge-output-format", "mkv",
		},
	},
	{
		Name:        "High Quality",
		Specs:       "1080p • H.264 + AAC • MP4",
		Quality:     "high",
		Description: "High Quality — 1080p • H.264 + AAC • MP4 (high)",
		Flags: []string{
			"-f", "bv*[height<=1080][ext=mp4]+ba[ext=m4a]/b[height<=1080][ext=mp4]",
			"--merge-output-format", "mp4",
		},
	},
	{
		Name:        "Balanced",
		Specs:       "720p • H.264 + AAC • MP4",
		Quality:     "good",
		Description: "Balanced — 720p • H.264 + AAC • MP4 (good)",
		Flags: []string{
			"-f", "bv*[height<=720][ext=mp4]+ba[ext=m4a]/b[height<=720][ext=mp4]",
			"--merge-output-format", "mp4",
		},
	},
	{
		Name:        "Mobile Saver",
		Specs:       "≤720p • H.264 • MP4 • low bitrate",
		Quality:     "mobile",
		Description: "Mobile Saver — ≤720p • H.264 • MP4 • low bitrate (mobile)",
		Flags: []string{
			"-f", "bv*[height<=720][vcodec^=avc1][ext=mp4]+ba[ext=m4a]/b[height<=720][ext=mp4]",
			"--merge-output-format", "mp4",
		},
	},
	{
		Name:        "Audio Only",
		Specs:       "Opus / MP3 • 160–320 kbps",
		Quality:     "audio",
		Description: "Audio Only — Opus / MP3 • 160–320 kbps (audio)",
		Flags: []string{
			"-f", "bestaudio[ext=m4a]/bestaudio[ext=mp3]/bestaudio",
			"-x", "--audio-format", "mp3",
			"--audio-quality", "0",
		},
	},
	{
		Name:        "Archive",
		Specs:       "Best video + best audio + metadata + thumbnails + subtitles",
		Quality:     "archive",
		Description: "Archive — Best video + best audio + metadata + thumbnails + subtitles (archive)",
		Flags: []string{
			"-f", "bestvideo+bestaudio/best",
			"--write-info-json",
			"--write-thumbnail",
			"--write-subs",
			"--write-auto-subs",
			"--merge-output-format", "mkv",
		},
	},
}

func GetProfiles() []Profile {
	return profiles
}

func GetDefaultProfile() *Profile {
	for i := range profiles {
		if profiles[i].Quality == "good" {
			return &profiles[i]
		}
	}
	if len(profiles) > 0 {
		return &profiles[0]
	}
	return nil
}

func GetProfileByQuality(quality string) *Profile {
	for i := range profiles {
		if profiles[i].Quality == quality {
			return &profiles[i]
		}
	}
	return GetDefaultProfile()
}
