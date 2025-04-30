package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/w0/ac/helpers"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		plistPath, err := cmd.Flags().GetString("plist")
		if err != nil {
			log.Fatalf("Failed to read flag plist: %v", err)
		}

		ac, err := helpers.ReadPlist(plistPath)

		log.Printf("ac: %d", len(ac.Packages))

		out, _ := cmd.Flags().GetString("output")

		log.Printf("output: %s", out)

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	downloadCmd.MarkFlagFilename("plist", "plist")
	downloadCmd.MarkFlagRequired("plist")

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get PWD: %v", err)
	}

	downloadCmd.PersistentFlags().StringP("output", "o", pwd, "Path to the location to download audio content")
	downloadCmd.MarkFlagDirname("output")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
