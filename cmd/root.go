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
		Example: `vip-video-downloader https://youku.com/v/xyz.html
vip-video-downloader https://youku.com/v/xyz.html -o /to/path -k -f mkv
vip-video-downloader https://youku.com/v/xyz.html -c=F
vip-video-downloader https://youku.com/v/xyz.html -F="/to/path/ffmpeg"`,
		Version: "1.0.0",
		Run:     run,
	}

	verbose             bool
	keepTS              bool
	outputDir           string
	convert             bool
	downloadChannel     string
	targetFormat        string
	ffmpegPath          string
	downloadConcurrency int
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "print verbose log")
	rootCmd.PersistentFlags().BoolVarP(&keepTS, "keep", "k", false, "keep TS files")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory")
	rootCmd.PersistentFlags().BoolVarP(&convert, "convert", "c", true, "convert downloaded video")
	rootCmd.PersistentFlags().StringVarP(&downloadChannel, "download-channel", "C", "lqiyi", "download video channel: lqiyi, ...")
	rootCmd.PersistentFlags().StringVarP(&targetFormat, "target-format", "f", "mp4", "convert to target format video: mp4, mkv, avi, ...")
	rootCmd.PersistentFlags().StringVarP(&ffmpegPath, "ffmpeg-path", "F", "ffmpeg", "The FFmpeg binary path, default auto detected in $PATH")
	rootCmd.PersistentFlags().IntVar(&downloadConcurrency, "download-concurrency", 25, "Download video concurrency")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(_ *cobra.Command, args []string) {
	if len(args) <= 0 {
		panic("error: require URL")
	}
	URL := args[0]
	channel := channel.GetChannel(downloadChannel)
	if channel == nil {
		panic("error: not support channel:" + downloadChannel)
	}
	parsedURL := channel.Parse(URL)
	tasker, err := download.NewTask(outputDir, parsedURL, verbose)
	if err != nil {
		panic(err)
	}
	err = tasker.Start(downloadConcurrency)
	if err != nil {
		panic(err)
	}
	outputFile := tasker.OutputFile()
	if convert {
		targetFilePath := filepath.Join(filepath.Dir(outputFile), "video."+strings.ToLower(targetFormat))
		ffmpegCmd := exec.Command(ffmpegPath, "-y", "-i", outputFile, "-c", "copy", targetFilePath)
		if verbose {
			ffmpegCmd.Stdout = os.Stdout
			ffmpegCmd.Stderr = os.Stdout
		}
		if err := ffmpegCmd.Run(); err != nil {
			panic(err)
		}
		if !keepTS {
			_ = os.RemoveAll(outputFile)
		}
		outputFile = targetFilePath
	}
	fmt.Println("[success] " + outputFile)
}
