package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w0/ac/helpers"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		plist := cmd.Flag("plist")
		ac, _ := helpers.ReadPlist(plist.Value.String())

		pkgName, _ := cmd.Flags().GetStringArray("name")

		fmt.Printf("%+v\n", ac.GetPackageByName(pkgName))

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content")
	listCmd.PersistentFlags().StringArrayP("name", "n", []string{}, "Package names to display")
	listCmd.PersistentFlags().StringArrayP("packageId", "i", []string{}, "Package ids to display")
	listCmd.PersistentFlags().BoolP("optional", "o", false, "Show only optional audio content")
	listCmd.PersistentFlags().BoolP("mandatory", "m", false, "Show only madatory audio content")
	listCmd.MarkFlagsMutuallyExclusive("name", "packageId")
}
