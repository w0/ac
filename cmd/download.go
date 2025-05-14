package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
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

		downloadLimit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			log.Fatalf("Failed to parse limit: %v", err)
		}

		pipeline := helpers.BuildFilterPipeline(cmd)

		filtered := pipeline(ac.Packages)

		downloadDir, _ := cmd.PersistentFlags().GetString("output")

		downloadContent(downloadDir, downloadLimit, filtered)

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

func downloadContent(output string, limit int, pkgs map[string]audiocontent.Packages) {

	log.Printf("Downloading %d packages.", len(pkgs))

	var wg sync.WaitGroup

	progress := mpb.New(mpb.WithAutoRefresh(), mpb.WithWaitGroup(&wg))

	wg.Add(len(pkgs))

	// TODO: fix goroutine downloads, ends before all downloads finish.
	//limitChan := make(chan struct{}, limit)

	for pkgName, values := range pkgs {

		resp, err := http.Get(string(values.DownloadName))
		if err != nil {
			log.Fatalf("Failed to fetch %s: %v", pkgName, err)
		}

		defer resp.Body.Close()

		bar := progress.New(
			resp.ContentLength,
			mpb.BarStyle().TipOnComplete(),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(pkgName),
				decor.Counters(decor.SizeB1024(0), " %.2f/%.2f"),
			),
			mpb.AppendDecorators(decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_MMSS, 0, decor.WCSyncWidth), "done!")),
		)

		func(bar *mpb.Bar) {
			wg.Done()

			proxy := bar.ProxyReader(resp.Body)

			defer proxy.Close()

			fileName := path.Base(string(values.DownloadName))

			outfile, err := os.Create(path.Join(output, fileName))
			if err != nil {
				log.Fatalf("Failed to create outfile %v", err)
			}

			defer outfile.Close()

			io.Copy(outfile, proxy)

		}(bar)

	}

}
