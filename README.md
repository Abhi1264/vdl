# VidForge — Universal Video Downloader

![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)
![GitHub Release](https://img.shields.io/github/v/release/Abhi1264/vidforge)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Homebrew](https://img.shields.io/homebrew/v/vidforge)

VidForge is a terminal-based video downloader application written in Go that provides an intuitive TUI (Terminal User Interface) for downloading videos from various platforms using `yt-dlp` as the backend. It features concurrent download management, quality profiles, SponsorBlock integration, and cross-platform support.

## Features

- 🎬 **Universal Video Downloading**: Supports all platforms that `yt-dlp` supports (YouTube, Vimeo, Twitter, TikTok, and 1000+ more)
- 🎨 **Beautiful TUI**: Modern terminal interface built with Bubbletea
- 📊 **Concurrent Downloads**: Download multiple videos simultaneously (up to 10 concurrent jobs)
- 🎚️ **Quality Profiles**: Pre-configured profiles for different use cases (Best Quality, High Quality, Balanced, Mobile Saver, Audio Only, Archive)
- 🚫 **SponsorBlock Integration**: Automatically removes sponsor segments from YouTube videos
- 🔄 **Resume Support**: Automatically resumes interrupted downloads
- 📦 **Self-Contained**: Includes embedded `yt-dlp` and `ffmpeg` binaries - no external dependencies required
- 🖥️ **Cross-Platform**: Works on macOS, Linux, and Windows

## Installation

### Homebrew (macOS/Linux)

```bash
brew install vidforge
```

### Download from GitHub Releases

Download the pre-built binary for your platform from the [Releases](https://github.com/Abhi1264/vidforge/releases) page:

- **macOS**: `vidforge_*_darwin_amd64.tar.gz` or `vidforge_*_darwin_arm64.tar.gz`
- **Linux**: `vidforge_*_linux_amd64.tar.gz` or `vidforge_*_linux_arm64.tar.gz`
- **Windows**: `vidforge_*_windows_amd64.tar.gz` or `vidforge_*_windows_arm64.tar.gz`

Extract and run:
```bash
tar -xzf vidforge_*_linux_amd64.tar.gz
./vidforge
```

### Build from Source

**Prerequisites:**
- Go 1.25.5 or later
- `curl` or `wget` for downloading embedded binaries

```bash
git clone https://github.com/Abhi1264/vidforge.git
cd vidforge

# Download embedded binaries (required before building)
bash scripts/download-binaries.sh

# Build the binary
go build -o vidforge ./cmd/vidforge
./vidforge
```

**Note**: The `download-binaries.sh` script downloads `yt-dlp` and `ffmpeg` binaries for all platforms, which are then embedded into the Go binary during compilation. This ensures the final binary is self-contained and doesn't require external dependencies.

### Dependencies

VidForge includes embedded binaries for `yt-dlp` and `ffmpeg`, so **no external dependencies are required**. The binaries are automatically extracted to a local `bin` directory on first run.

**How it works:**
1. On first launch, VidForge checks if `yt-dlp` and `ffmpeg` are available in your system PATH
2. If not found, it extracts the embedded binaries to `{app_directory}/bin/`
3. The extracted binaries are used for all download operations

**Fallback behavior:**
- If you have system-installed versions of `yt-dlp` or `ffmpeg` in your PATH, those will be used instead
- This allows you to use updated versions if you prefer

## Usage

### Starting VidForge

Simply run the executable:

```bash
./vidforge
```

### Basic Workflow

1. **Enter a URL**: Type or paste a video URL in the input field
2. **Press Enter**: Start the download
3. **Monitor Progress**: Watch real-time progress bars for each download
4. **Manage Downloads**: Navigate between downloads, pause/cancel, or change settings

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Start download with the entered URL |
| `↑` / `↓` or `j` / `k` | Navigate between download jobs |
| `p` | Pause/cancel the selected download |
| `s` | Toggle SponsorBlock (YouTube only) |
| `f` | Open profile selection menu |
| `?` | Toggle help screen |
| `q` / `Ctrl+C` | Quit the application |

### Download Profiles

VidForge includes 6 pre-configured download profiles optimized for different use cases:

1. **Best Quality** (`best`)
   - 4K/8K resolution
   - VP9/AV1 video codec + Opus audio
   - MKV container format
   - Best available quality

2. **High Quality** (`high`)
   - 1080p resolution
   - H.264 video codec + AAC audio
   - MP4 container format
   - Good balance of quality and compatibility

3. **Balanced** (`good`) - Default
   - 720p resolution
   - H.264 video codec + AAC audio
   - MP4 container format
   - Optimal for most use cases

4. **Mobile Saver** (`mobile`)
   - ≤720p resolution
   - H.264 video codec
   - MP4 container format
   - Low bitrate for mobile devices

5. **Audio Only** (`audio`)
   - Opus/MP3 format
   - 160-320 kbps bitrate
   - Extracts audio only

6. **Archive** (`archive`)
   - Best video + best audio
   - Includes metadata, thumbnails, and subtitles
   - MKV container format
   - Complete archival package

To change profiles:
- Press `f` to open the profile selection menu
- Use `↑`/`↓` to navigate or press `1-6` to select directly
- Press `Enter` to confirm

### SponsorBlock

VidForge includes SponsorBlock integration for YouTube videos. When enabled (default), it automatically removes:
- Sponsor segments
- Self-promotion segments
- Interaction requests
- Intro segments
- Outro segments

Toggle with the `s` key. SponsorBlock is automatically enabled for YouTube URLs and disabled for other platforms.

## Architecture

### Project Structure

```
vidforge/
├── cmd/
│   └── vidforge/
│       └── main.go          # Application entry point
├── internal/
│   ├── bootstrap/
│   │   ├── deps.go          # Dependency installation logic
│   │   └── os.go            # OS detection utilities
│   ├── downloader/
│   │   ├── formats.go       # Format listing functionality
│   │   ├── job.go           # Download job execution
│   │   ├── manager.go       # Concurrent download manager
│   │   ├── profile.go       # Quality profiles
│   │   └── url.go           # URL validation utilities
│   └── ui/
│       ├── formats.go       # Format list item implementation
│       ├── help.go          # Help text
│       ├── model.go         # Bubbletea model (state)
│       ├── styles.go        # UI styling
│       ├── update.go        # Event handling (update logic)
│       └── view.go          # UI rendering
├── go.mod
├── go.sum
└── README.md
```

### Core Components

#### 1. Bootstrap Package (`internal/bootstrap`)

Handles dependency management and binary extraction:

- **`deps.go`**: Manages embedded binaries and extraction
  - Embeds `yt-dlp` and `ffmpeg` binaries for all platforms
  - Extracts binaries to local directory on first run
  - Checks system PATH first, falls back to embedded binaries
  - Platform-specific binary selection based on OS and architecture

- **`os.go`**: OS and architecture detection utilities
  - `DetectOS()`: Returns current OS and architecture
  - `IsMacOS()`, `IsLinux()`, `IsWindows()`: Platform checks
  - `IsIntel()`, `IsARM()`: Architecture checks

#### 2. Downloader Package (`internal/downloader`)

Core download functionality:

- **`manager.go`**: Concurrent download manager
  - Manages a queue of download jobs
  - Supports multiple concurrent downloads (configurable workers)
  - Provides progress updates via channels
  - Supports job cancellation

- **`job.go`**: Individual download job execution
  - Runs `yt-dlp` processes
  - Parses progress output using regex
  - Handles context cancellation
  - Supports resume functionality

- **`profile.go`**: Download quality profiles
  - Pre-defined quality profiles with yt-dlp flags
  - Profile selection by quality string
  - Default profile selection

- **`formats.go`**: Format listing utility
  - Parses `yt-dlp -F` output
  - Returns available formats for a URL (not currently used in UI)

- **`url.go`**: URL utilities
  - YouTube URL detection
  - Platform-specific handling

#### 3. UI Package (`internal/ui`)

Terminal user interface built with Bubbletea:

- **`model.go`**: Application state
  - Text input for URLs
  - Job tracking (map of job states)
  - Profile selection state
  - SponsorBlock toggle state
  - Help/profile menu visibility

- **`update.go`**: Event handling
  - Keyboard input processing
  - Progress updates from download manager
  - State transitions
  - Job submission logic

- **`view.go`**: UI rendering
  - Main view rendering
  - Progress bar rendering
  - Profile selection view
  - Help screen rendering
  - Job list rendering

- **`styles.go`**: UI styling
  - Lipgloss styles for different UI elements
  - Color schemes and formatting

- **`help.go`**: Help text content

- **`formats.go`**: Format list item type definition for Bubbletea lists (utility type)

### Data Flow

1. **Initialization**:
   - `main.go` ensures dependencies are available via `bootstrap.Ensure()`
   - Extracts embedded binaries if system versions aren't found
   - Creates a new Bubbletea program with `ui.NewModel()`
   - Initializes download manager with 10 worker goroutines

2. **Download Submission**:
   - User enters URL and presses Enter
   - Model validates URL and creates a new job
   - Job is submitted to the download manager queue
   - Manager assigns job to an available worker

3. **Download Execution**:
   - Worker goroutine uses `bootstrap.GetCommandPath()` to locate `yt-dlp`
   - Executes `yt-dlp` command (from PATH or extracted binary) with appropriate flags
   - Progress is parsed from stdout using regex
   - Progress updates are sent via channel to the UI model

4. **UI Updates**:
   - UI listens to progress channel
   - Model updates job states on progress messages
   - View re-renders with updated progress bars and status

### Concurrency Model

VidForge uses a worker pool pattern for concurrent downloads:

- **Download Manager**: Maintains a channel-based queue of jobs
- **Worker Goroutines**: 10 workers process jobs concurrently
- **Progress Channel**: Single channel for all progress updates
- **Context Cancellation**: Each job has a context for cancellation support

## Configuration

VidForge uses sensible defaults and doesn't require configuration files:

- **Concurrent Downloads**: 10 workers (configurable in `ui/model.go`)
- **Default Profile**: Balanced (720p, H.264, MP4)
- **SponsorBlock**: Enabled by default for YouTube URLs
- **Resume**: Enabled by default for all downloads

## Development

### Building

```bash
go build -o vidforge ./cmd/vidforge
```

### Testing

Run the application to test:

```bash
go run ./cmd/vidforge
```

### Adding New Profiles

Edit `internal/downloader/profile.go` and add a new `Profile` to the `profiles` slice:

```go
{
    Name:        "Custom Profile",
    Specs:       "Your specs description",
    Quality:     "custom",
    Description: "Full description",
    Flags: []string{
        "-f", "your-format-string",
        "--merge-output-format", "mp4",
    },
}
```

### Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/bubbles`: UI components (textinput)
- `github.com/charmbracelet/lipgloss`: Terminal styling

Embedded dependencies (included in binary):
- `yt-dlp`: Video downloader (embedded for all platforms)
- `ffmpeg`: Media processing (embedded for all platforms)

These are automatically extracted on first run if not found in system PATH.

## Limitations

1. **Format Selection**: Currently uses profiles only; manual format selection not available
3. **Error Handling**: Basic error display; detailed error recovery may need improvement
4. **Concurrent Workers**: Fixed at 10 workers (configurable in code)
5. **Binary Size**: The final binary is larger (~50-100MB) due to embedded dependencies

## Future Enhancements

Potential improvements:

- [ ] Configurable download directory
- [ ] Custom format selection UI
- [ ] Download history persistence
- [ ] Playlist support with queue management
- [ ] Subtitle selection and management
- [ ] Config file for user preferences
- [ ] Better error messages and recovery
- [ ] Download speed indicators
- [ ] ETA calculation and display
- [ ] Optional: Use system-installed binaries by default (flag to prefer embedded)

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Powered by [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- Uses [SponsorBlock](https://sponsor.ajay.app/) for YouTube ad removal

