package cmd

import (
	"fmt"
	"github.com/billcoding/vip-video-downloader/channel"
	m3u8DL "github.com/billcoding/vip-video-downloader/m3u8/download"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	downloadCmd = &cobra.Command{
		Use:     "download URL",
		Aliases: []string{"d"},
		Short:   "Download Vip Video or M3u8",
		Long: `Vip Video Download
Supports: YouKu | QQ | SoHu TV | Mango TV ...
If you need more, please contact me.`,
		Example: `
- normal: vip-video-downloader download "https://abc.com/v/xyz.html"
- output dir: vip-video-downloader download -d "/to/path" "https://abc.com/v/xyz.html"
- output file: vip-video-downloader download -o "my-video" "https://abc.com/v/xyz.html"
- no convert(.ts): vip-video-downloader download -c=F "https://abc.com/v/xyz.html"
- convert to mkv format: vip-video-downloader download -f="mkv" "https://abc.com/v/xyz.html"
- special FFmpeg path: vip-video-downloader download -F="/to/path/ffmpeg" "https://abc.com/v/xyz.html"
- m3u8 URL: vip-video-downloader download -m "https://abc.com/index.m3u8"
- m3u8 File: vip-video-downloader download -m -M "/to/path/v.m3u8"`,
		Run: downloadRun,
	}

	verbose    bool
	keepTS     bool
	tsDir      string
	outputDir  string
	outputFile string

	m3u8     bool
	m3u8File bool

	downloadChannel     string
	downloadConcurrency int

	convert       bool
	convertFormat string

	ffmpegPath string
)

func init() {
	downloadCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Print verbose log")
	downloadCmd.PersistentFlags().BoolVarP(&keepTS, "keep", "k", false, "Keep TS files")
	downloadCmd.PersistentFlags().StringVar(&tsDir, "ts-dir", "ts", "TS files directory")
	downloadCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "d", "", "Output directory")
	downloadCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "video", "Output file name without extension")

	downloadCmd.PersistentFlags().BoolVarP(&m3u8, "m3u8", "m", false, "m3u8 direct")
	downloadCmd.PersistentFlags().BoolVarP(&m3u8File, "m3u8-file", "M", false, "m3u8 file")

	downloadCmd.PersistentFlags().StringVarP(&downloadChannel, "download-channel", "C", "c1", "Download video channel: c1, ...")
	downloadCmd.PersistentFlags().IntVarP(&downloadConcurrency, "download-concurrency", "N", 25, "Download video concurrency")

	downloadCmd.PersistentFlags().BoolVarP(&convert, "convert", "c", true, "Convert downloaded video")
	downloadCmd.PersistentFlags().StringVarP(&convertFormat, "convert-format", "f", "mp4", "Convert to target format video: mp4, mkv, avi, ...")

	downloadCmd.PersistentFlags().StringVarP(&ffmpegPath, "ffmpeg-path", "F", "ffmpeg", "The FFmpeg binary path, default auto detected in $PATH")
}

func downloadRun(_ *cobra.Command, args []string) {
	if len(args) <= 0 {
		panic("error: require URL")
	}
	URL := args[0]
	resultFile := ""
	if !m3u8 {
		if c := channel.GetChannel(downloadChannel); c == nil {
			panic("error: not support channel:" + downloadChannel)
		} else {
			URL = c.Parse(URL)
			if verbose {
				fmt.Printf("[m3u8]%s\n", URL)
			}
		}
	}
	tasker, err := m3u8DL.NewTask(outputDir, outputFile, URL, tsDir, m3u8File, verbose)
	if err != nil {
		panic(err)
	}
	err = tasker.Start(downloadConcurrency)
	if err != nil {
		panic(err)
	}
	resultFile = tasker.OutputFile()
	targetFilePath := filepath.Join(filepath.Dir(resultFile), outputFile+"."+strings.ToLower(convertFormat))
	if convert && !strings.EqualFold(resultFile, targetFilePath) {
		ffmpegCmd := exec.Command(ffmpegPath, "-y", "-i", resultFile, "-c", "copy", targetFilePath)
		if verbose {
			ffmpegCmd.Stdout = os.Stdout
			ffmpegCmd.Stderr = os.Stdout
		}
		if err = ffmpegCmd.Run(); err != nil {
			panic(err)
		}
		if !keepTS {
			_ = os.RemoveAll(resultFile)
		}
		resultFile = targetFilePath
	}
	fmt.Println("[success] " + resultFile)
}
