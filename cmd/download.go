package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/w0/ac/helpers"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download available audio content pkgs.",
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

		progress := mpb.NewWithContext(
			cmd.Context(),
			mpb.WithRefreshRate(240*time.Millisecond),
			mpb.WithWidth(60))

		_, err = helpers.DownloadPackages(progress, &filtered, downloadDir)

		progress.Wait()

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

}
