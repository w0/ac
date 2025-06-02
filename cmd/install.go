package cmd

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/w0/ac/helpers"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install available audio content pkgs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		plistPath, err := cmd.Flags().GetString("plist")
		if err != nil {
			log.Fatalf("Failed to read flag plist: %s", err)
		}

		ac, err := helpers.ReadPlist(plistPath)
		if err != nil {
			log.Fatalf("Failed to parse audio content: %s", err)
		}

		pipeline := helpers.BuildFilterPipeline(cmd)

		filtered := pipeline(ac.Packages)

		downloadDir, _ := cmd.PersistentFlags().GetString("output")

		downProgress := mpb.NewWithContext(
			cmd.Context(),
			mpb.WithRefreshRate(240*time.Millisecond),
			mpb.WithWidth(60))

		installers, err := helpers.DownloadPackages(downProgress, &filtered, downloadDir)
		if err != nil {
			log.Fatalf("Failed download pkgs: %s", err)
		}

		downProgress.Wait()

		installProgress := mpb.NewWithContext(
			cmd.Context(),
			mpb.WithRefreshRate(60*time.Millisecond),
			mpb.WithWidth(60))

		err = helpers.InstallPackages(installProgress, installers)
		if err != nil {
			log.Fatalf("Failed installing pkgs: %s", err)
		}

		installProgress.Wait()

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	if os := runtime.GOOS; os != "darwin" {
		log.Fatalf("install only supported on macOS.\n")
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get PWD: %v", err)
	}

	installCmd.PersistentFlags().StringP("output", "d", pwd, "Path to the location to download audio content")
	installCmd.MarkPersistentFlagDirname("output")

	installCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	installCmd.MarkPersistentFlagFilename("plist", "plist")
	installCmd.MarkPersistentFlagRequired("plist")

	installCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Package names to download")
	installCmd.PersistentFlags().StringSliceP("packageId", "i", []string{}, "Package ids to download")
	installCmd.MarkFlagsMutuallyExclusive("name", "packageId")

	installCmd.PersistentFlags().BoolP("optional", "o", false, "Download only optional audio content")
	installCmd.PersistentFlags().BoolP("mandatory", "m", false, "Download only madatory audio content")
	installCmd.MarkFlagsMutuallyExclusive("optional", "mandatory")

}
