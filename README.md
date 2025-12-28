# VidForge — Universal Video Downloader

VidForge is a terminal-based video downloader application written in Go that provides an intuitive TUI (Terminal User Interface) for downloading videos from various platforms using `yt-dlp` as the backend. It features concurrent download management, quality profiles, SponsorBlock integration, and cross-platform support.

## Features

- 🎬 **Universal Video Downloading**: Supports all platforms that `yt-dlp` supports (YouTube, Vimeo, Twitter, TikTok, and 1000+ more)
- 🎨 **Beautiful TUI**: Modern terminal interface built with Bubbletea
- 📊 **Concurrent Downloads**: Download multiple videos simultaneously (up to 3 concurrent jobs)
- 🎚️ **Quality Profiles**: Pre-configured profiles for different use cases (Best Quality, High Quality, Balanced, Mobile Saver, Audio Only, Archive)
- 🚫 **SponsorBlock Integration**: Automatically removes sponsor segments from YouTube videos
- 🔄 **Resume Support**: Automatically resumes interrupted downloads
- 🔧 **Auto-Dependency Management**: Automatically installs `yt-dlp` and `ffmpeg` if not present
- 🖥️ **Cross-Platform**: Works on macOS, Linux, and Windows

## Installation

### Prerequisites

- Go 1.25.5 or later
- Homebrew (macOS), apt-get/dnf/pacman (Linux), or PowerShell (Windows) for automatic dependency installation

### Build from Source

```bash
git clone https://github.com/Abhi1264/vidforge.git
cd vidforge
go build -o vidforge ./cmd/vidforge
./vidforge
```

### Dependencies

VidForge automatically checks for and installs the following dependencies:

- **yt-dlp**: Video downloader backend
- **ffmpeg**: Media processing tool

If these are not found, VidForge will attempt to install them:
- **macOS**: Uses Homebrew (`brew install yt-dlp ffmpeg`)
- **Linux**: Uses apt-get, dnf, or pacman depending on your distribution
- **Windows**: Downloads yt-dlp.exe to the application directory (ffmpeg Windows installation not fully implemented)

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

Handles dependency management and OS detection:

- **`deps.go`**: Ensures `yt-dlp` and `ffmpeg` are installed
  - Cross-platform installation logic
  - Package manager detection (Homebrew, apt-get, dnf, pacman)
  - Automatic download for Windows

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

- **`formats.go`**: Format listing
  - Parses `yt-dlp -F` output
  - Returns available formats for a URL

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

- **`formats.go`**: Format list item implementation for Bubbletea lists

### Data Flow

1. **Initialization**:
   - `main.go` checks for dependencies via `bootstrap.Ensure()`
   - Creates a new Bubbletea program with `ui.NewModel()`
   - Initializes download manager with 3 worker goroutines

2. **Download Submission**:
   - User enters URL and presses Enter
   - Model validates URL and creates a new job
   - Job is submitted to the download manager queue
   - Manager assigns job to an available worker

3. **Download Execution**:
   - Worker goroutine executes `yt-dlp` command with appropriate flags
   - Progress is parsed from stdout using regex
   - Progress updates are sent via channel to the UI model

4. **UI Updates**:
   - UI listens to progress channel
   - Model updates job states on progress messages
   - View re-renders with updated progress bars and status

### Concurrency Model

VidForge uses a worker pool pattern for concurrent downloads:

- **Download Manager**: Maintains a channel-based queue of jobs
- **Worker Goroutines**: 3 workers process jobs concurrently
- **Progress Channel**: Single channel for all progress updates
- **Context Cancellation**: Each job has a context for cancellation support

## Configuration

VidForge uses sensible defaults and doesn't require configuration files:

- **Concurrent Downloads**: 3 workers (configurable in `ui/model.go`)
- **Default Profile**: Balanced (720p, H.264, MP4)
- **SponsorBlock**: Enabled by default for YouTube URLs
- **Resume**: Enabled by default for all downloads

## Development

### Building

```bash
go build -o vidforge ./cmd/vidforge
```

### Testing

While no formal tests are included, you can test individual components:

```bash
# Test dependency detection
go run ./cmd/vidforge  # Will check for yt-dlp and ffmpeg

# Test format listing (requires yt-dlp)
go run -c 'import "github.com/Abhi1264/vidforge/internal/downloader"; ...'
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

External dependencies (installed automatically):
- `yt-dlp`: Video downloader
- `ffmpeg`: Media processing

## Limitations

1. **Windows Support**: FFmpeg Windows installation is not fully implemented (downloads yt-dlp.exe only)
2. **Format Selection**: Currently uses profiles only; manual format selection not available
3. **Download Location**: Uses yt-dlp default download location (current directory)
4. **Error Handling**: Basic error display; detailed error recovery may need improvement
5. **Concurrent Workers**: Fixed at 3 workers (configurable in code)

## Future Enhancements

Potential improvements:

- [ ] Configurable download directory
- [ ] Custom format selection UI
- [ ] Download history persistence
- [ ] Playlist support with queue management
- [ ] Subtitle selection and management
- [ ] Config file for user preferences
- [ ] Windows FFmpeg installation completion
- [ ] Better error messages and recovery
- [ ] Download speed indicators
- [ ] ETA calculation and display

## License

This project appears to be open source. Please check with the repository owner for specific licensing information.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Powered by [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- Uses [SponsorBlock](https://sponsor.ajay.app/) for YouTube ad removal

