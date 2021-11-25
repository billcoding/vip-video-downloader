package cmd

import (
	"fmt"
	"github.com/billcoding/vip-video-downloader/channel"
	"github.com/billcoding/vip-video-downloader/download"
	m3u8DL "github.com/billcoding/vip-video-downloader/m3u8/download"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	downloadCmd = &cobra.Command{
		Use:     "download",
		Aliases: []string{"d"},
		Short:   "Download Vip Video or M3u8",
		Long: `Vip Video Download
Supports: YouKu | QQ | SoHu TV | Mango TV ...
If you need more, please contact me.`,
		Example: `normal: vip-video-downloader download https://youku.com/v/xyz.html
output dir: vip-video-downloader download https://youku.com/v/xyz.html -d /to/path
output file: vip-video-downloader download https://youku.com/v/xyz.html -o my-video
no convert: vip-video-downloader download https://youku.com/v/xyz.html -c=F
convert to mkv format: vip-video-downloader download https://youku.com/v/xyz.html -f="mkv"
special FFmpeg path: vip-video-downloader download https://youku.com/v/xyz.html -F="/to/path/ffmpeg"
no use download channel: vip-video-downloader download https://example.com/index.m3u8 -U=F
use download channel: vip-video-downloader download https://example.com/index.m3u8 -C=lqiyi -N=100
`,
		Run: downloadRun,
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
	downloadCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Print verbose log")
	downloadCmd.PersistentFlags().BoolVarP(&keepTS, "keep", "k", false, "Keep TS files")
	downloadCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "d", "", "Output directory")
	downloadCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "video", "Output file name without extension")

	downloadCmd.PersistentFlags().BoolVarP(&useDownloadChannel, "use-download-channel", "U", true, "Use download channel")
	downloadCmd.PersistentFlags().StringVarP(&downloadChannel, "download-channel", "C", "lqiyi", "Download video channel: lqiyi, ...")
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
	URLExt := ""
	isM3u8 := false
	resultFile := ""
	if useDownloadChannel {
		if c := channel.GetChannel(downloadChannel); c == nil {
			panic("error: not support channel:" + downloadChannel)
		} else {
			URL, URLExt, isM3u8 = c.Parse(URL)
		}
	}

	if isM3u8 {
		tasker, err := m3u8DL.NewTask(outputDir, outputFile, URL, verbose)
		if err != nil {
			panic(err)
		}
		err = tasker.Start(downloadConcurrency)
		if err != nil {
			panic(err)
		}
		resultFile = tasker.OutputFile()
	} else {
		resultFile = filepath.Join(outputDir, outputFile+"."+URLExt)
		dlER := download.Downloader{
			URL:    URL,
			Output: resultFile,
		}
		dlER.Start(verbose)
	}

	var err error
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
