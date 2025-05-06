package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/w0/ac/audiocontent"
	"github.com/w0/ac/helpers"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		plistPath, err := cmd.Flags().GetString("plist")
		if err != nil {
			log.Fatalf("Failed to read flag plist: %v", err)
		}

		ac, err := helpers.ReadPlist(plistPath)
		if err != nil {
			log.Fatalf("Failed to parse audio content: %v", err)
		}

		pipeline := helpers.BuildFilterPipeline(cmd)

		filtered := pipeline(ac.Packages)

		downloadDir, _ := cmd.PersistentFlags().GetString("output")
		downloadLimit, _ := cmd.PersistentFlags().GetInt("limit")

		jobs := make(chan audiocontent.Packages, downloadLimit)

		for w := 1; w <= 3; w++ {
			go downloadContent(downloadDir, jobs)
		}

		for _, v := range filtered {
			jobs <- v
		}
		close(jobs)

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get PWD: %v", err)
	}

	downloadCmd.PersistentFlags().StringP("output", "d", pwd, "Path to the location to download audio content")
	downloadCmd.MarkPersistentFlagDirname("output")

	downloadCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	downloadCmd.MarkPersistentFlagFilename("plist", "plist")
	downloadCmd.MarkPersistentFlagRequired("plist")

	downloadCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Package names to download")
	downloadCmd.PersistentFlags().StringSliceP("packageId", "i", []string{}, "Package ids to download")
	downloadCmd.MarkFlagsMutuallyExclusive("name", "packageId")

	downloadCmd.PersistentFlags().BoolP("optional", "o", false, "Download only optional audio content")
	downloadCmd.PersistentFlags().BoolP("mandatory", "m", false, "Download only madatory audio content")
	downloadCmd.MarkFlagsMutuallyExclusive("optional", "mandatory")

	downloadCmd.PersistentFlags().IntP("limit", "l", 3, "Limit the concurrent download of audio content")
}

func downloadContent(output string, jobs <-chan audiocontent.Packages) {
	for i := range jobs {
		file := path.Join(output, path.Base(string(i.DownloadName)))

		f, _ := os.Create(file)
		defer f.Close()

		req, err := http.Get(string(i.DownloadName))
		if err != nil {
			log.Printf("FAILED: %s", i.DownloadName)
		}

		defer req.Body.Close()

		io.Copy(f, req.Body)
	}
}
