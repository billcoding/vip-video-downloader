package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "vip-video-downloader",
		Short:   "Vip Video Downloader",
		Long:    `Vip Video Downloader`,
		Version: "1.0.4",
	}
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(mergeCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
