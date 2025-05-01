package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

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

		pipeline := buildFilterPipeline(cmd)

		filtered := pipeline(ac.Packages)

		outputResults(cmd, filtered)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("plist", "p", "", "Path to plist containing audio content (required)")
	listCmd.MarkFlagRequired("plist")

	listCmd.PersistentFlags().StringSliceP("name", "n", []string{}, "Package names to display")
	listCmd.PersistentFlags().StringSliceP("packageId", "i", []string{}, "Package ids to display")
	listCmd.MarkFlagsMutuallyExclusive("name", "packageId")

	listCmd.PersistentFlags().BoolP("optional", "o", false, "Show only optional audio content")
	listCmd.PersistentFlags().BoolP("mandatory", "m", false, "Show only madatory audio content")
	listCmd.MarkFlagsMutuallyExclusive("optional", "mandatory")

	listCmd.PersistentFlags().BoolP("json", "j", false, "Output audio content info as json")
}

type Filter func(map[string]audiocontent.Packages) map[string]audiocontent.Packages

func buildFilterPipeline(cmd *cobra.Command) Filter {
	var filters []Filter

	if names, _ := cmd.PersistentFlags().GetStringSlice("name"); len(names) > 0 {
		filters = append(filters, func(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
			return filterByName(pkgs, names)
		})
	}

	if ids, _ := cmd.PersistentFlags().GetStringSlice("packageId"); len(ids) > 0 {
		filters = append(filters, func(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
			return filterById(pkgs, ids)
		})
	}

	if optional, _ := cmd.PersistentFlags().GetBool("optional"); optional {
		filters = append(filters, func(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
			return filterOptional(pkgs)
		})
	}

	if mandatory, _ := cmd.PersistentFlags().GetBool("mandatory"); mandatory {
		filters = append(filters, func(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
			return filterMandatory(pkgs)
		})
	}

	return func(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
		result := pkgs
		for _, filter := range filters {
			result = filter(result)
		}

		return result
	}

}

func filterByName(pkgs map[string]audiocontent.Packages, names []string) map[string]audiocontent.Packages {
	if len(names) == 0 {
		return pkgs
	}

	result := make(map[string]audiocontent.Packages)

	for _, name := range names {
		if val, ok := pkgs[name]; ok {
			result[name] = val
		}
	}

	return result
}

func filterById(pkgs map[string]audiocontent.Packages, ids []string) map[string]audiocontent.Packages {
	if len(ids) == 0 {
		return pkgs
	}

	result := make(map[string]audiocontent.Packages)

	for k, v := range pkgs {
		if slices.Contains(ids, v.PackageID) {
			result[k] = v
		}
	}

	return result
}

func filterOptional(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
	result := make(map[string]audiocontent.Packages)

	for k, v := range pkgs {
		if v.IsMandatory == false {
			result[k] = v
		}
	}

	return result
}

func filterMandatory(pkgs map[string]audiocontent.Packages) map[string]audiocontent.Packages {
	result := make(map[string]audiocontent.Packages)

	for k, v := range pkgs {
		if v.IsMandatory == true {
			result[k] = v
		}
	}

	return result
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
