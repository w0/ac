package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/w0/ac/audiocontent"
	"github.com/w0/ac/helpers"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list and search for available audio content",
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

		outputResults(cmd, filtered)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	listCmd.MarkPersistentFlagFilename("plist", "plist")
	listCmd.MarkPersistentFlagRequired("plist")

	listCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Package names to display")
	listCmd.PersistentFlags().StringSliceP("packageId", "i", []string{}, "Package ids to display")
	listCmd.MarkFlagsMutuallyExclusive("name", "packageId")

	listCmd.PersistentFlags().BoolP("optional", "o", false, "Show only optional audio content")
	listCmd.PersistentFlags().BoolP("mandatory", "m", false, "Show only madatory audio content")
	listCmd.MarkFlagsMutuallyExclusive("optional", "mandatory")

	listCmd.PersistentFlags().BoolP("json", "j", false, "Output audio content info as json")
}

func outputResults(cmd *cobra.Command, pkgs map[string]audiocontent.Packages) {
	outJson, _ := cmd.Flags().GetBool("json")

	if outJson {
		jsonOut, err := json.MarshalIndent(&pkgs, "", "  ")
		if err != nil {
			log.Fatalf("Failed to generate json: %v", err)
		}

		fmt.Println(string(jsonOut))
	} else {
		for k, v := range pkgs {
			fmt.Printf("\n-- %s --\n", k)
			fmt.Printf("  - ContainsAppleLoops: %t\n", v.ContainsAppleLoops)
			fmt.Printf("  - ContainsGarageBandLegacyInstruments: %t\n", v.ContainsGarageBandLegacyInstruments)
			fmt.Printf("  - DownloadName: %s\n", v.DownloadName)
			fmt.Printf("  - DownloadSize: %0.2f Mb\n", helpers.ConvertToMB(v.DownloadSize))
			fmt.Printf("  - FileCheck:\n")

			for _, fc := range v.FileCheck {
				fmt.Printf("    * %s\n", fc)
			}

			fmt.Printf("  - InstalledSize: %0.2f Mb\n", helpers.ConvertToMB(v.InstalledSize))
			fmt.Printf("  - IsMandatory: %t\n", v.IsMandatory)
			fmt.Printf("  - PackageID: %s\n", v.PackageID)
		}
	}
}
