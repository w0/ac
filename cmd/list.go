package cmd

import (
	"log"

	"github.com/spf13/cobra"
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

		pkgName, err := cmd.Flags().GetStringSlice("name")

		if err != nil {
			log.Fatalf("Failed to parse name: %v", err)
		}

		helpers.PrettyPrint(pkgName, ac)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	listCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Package names to display")
	listCmd.PersistentFlags().StringSliceP("packageId", "i", []string{}, "Package ids to display")
	listCmd.PersistentFlags().BoolP("optional", "o", false, "Show only optional audio content")
	listCmd.PersistentFlags().BoolP("mandatory", "m", false, "Show only madatory audio content")
	listCmd.MarkFlagsMutuallyExclusive("name", "packageId")
	listCmd.MarkFlagRequired("plist")
}
