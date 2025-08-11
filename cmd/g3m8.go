package cmd

import (
	"fmt"
	"g3u8/download"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "g3m8",
}

var downloadCmd = &cobra.Command{
	Use:   "download <url> ",
	Short: "download *.ts file with m3u8 url",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "url is required")
			os.Exit(1)
		}
		url := args[0]
		out, _ := cmd.Flags().GetString("out")
		baseUrl, _ := cmd.Flags().GetString("baseUrl")

		if out == "" {
			out = "./out.mp4"
		}

		err := download.M3u8Download(download.M3u8DownloadOptions{
			Out:     out,
			Url:     url,
			BaseUrl: baseUrl,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {

	downloadCmd.Flags().StringP("out", "o", "", "output file")
	downloadCmd.Flags().StringP("baseUrl", "u", "", "base url for playlist")
	rootCmd.AddCommand(downloadCmd)

}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
