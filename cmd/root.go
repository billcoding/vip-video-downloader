package cmd

import (
	"fmt"
	"github.com/billcoding/vip-video-downloader/channel"
	"github.com/billcoding/vip-video-downloader/m3u8/download"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Short: "Vip Video Downloader",
		Long: `Vip Video Downloader
Supports: YouKu | QQ | SoHu TV | Mango TV ...
If you need more, please contact me.`,
		Use: "vip-video-downloader URL",
		Example: `normal: vip-video-downloader https://youku.com/v/xyz.html
output dir: vip-video-downloader https://youku.com/v/xyz.html -d /to/path
output file: vip-video-downloader https://youku.com/v/xyz.html -o my-video
no convert: vip-video-downloader https://youku.com/v/xyz.html -c=F
convert to mkv format: vip-video-downloader https://youku.com/v/xyz.html -f="mkv"
special FFmpeg path: vip-video-downloader https://youku.com/v/xyz.html -F="/to/path/ffmpeg"
no use download channel: vip-video-downloader https://example.com/index.m3u8 -U=F
use download channel: vip-video-downloader https://example.com/index.m3u8 -C=lqiyi -N=100
`,
		Version: "1.0.0",
		Run:     run,
	}

	verbose    bool
	keepTS     bool
	outputDir  string
	outputFile string

	useDownloadChannel  bool
	downloadChannel     string
	downloadConcurrency int

	convert       bool
	convertFormat string

	ffmpegPath string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "print verbose log")
	rootCmd.PersistentFlags().BoolVarP(&keepTS, "keep", "k", false, "keep TS files")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "d", "", "output directory")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "video", "output file name without extension")

	rootCmd.PersistentFlags().BoolVarP(&useDownloadChannel, "use-download-channel", "U", true, "use download channel")
	rootCmd.PersistentFlags().StringVarP(&downloadChannel, "download-channel", "C", "lqiyi", "download video channel: lqiyi, ...")
	rootCmd.PersistentFlags().IntVarP(&downloadConcurrency, "download-concurrency", "N", 25, "download video concurrency")

	rootCmd.PersistentFlags().BoolVarP(&convert, "convert", "c", true, "convert downloaded video")
	rootCmd.PersistentFlags().StringVarP(&convertFormat, "convert-format", "f", "mp4", "convert to target format video: mp4, mkv, avi, ...")

	rootCmd.PersistentFlags().StringVarP(&ffmpegPath, "ffmpeg-path", "F", "ffmpeg", "The FFmpeg binary path, default auto detected in $PATH")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(_ *cobra.Command, args []string) {
	if len(args) <= 0 {
		panic("error: require URL")
	}
	URL := args[0]
	if useDownloadChannel {
		if c := channel.GetChannel(downloadChannel); c == nil {
			panic("error: not support channel:" + downloadChannel)
		} else {
			URL = c.Parse(URL)
		}
	}
	tasker, err := download.NewTask(outputDir, outputFile, URL, verbose)
	if err != nil {
		panic(err)
	}
	err = tasker.Start(downloadConcurrency)
	if err != nil {
		panic(err)
	}
	resultFile := tasker.OutputFile()
	if convert {
		targetFilePath := filepath.Join(filepath.Dir(resultFile), outputFile+"."+strings.ToLower(convertFormat))
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
