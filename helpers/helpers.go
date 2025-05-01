package helpers

import (
	"os"
	"slices"

	"github.com/micromdm/plist"
	"github.com/spf13/cobra"
	"github.com/w0/ac/audiocontent"
)

func ReadPlist(path string) (audiocontent.Content, error) {
	file, err := os.Open(path)
	if err != nil {
		return audiocontent.Content{}, err
	}

	defer file.Close()

	decoder := plist.NewDecoder(file)
	var ac audiocontent.Content

	err = decoder.Decode(&ac)
	if err != nil {
		return audiocontent.Content{}, err
	}

	return ac, nil
}

func ConvertToMB(n audiocontent.ContentSize) float64 {
	return float64(n) / (1 << 20)
}

type Filter func(map[string]audiocontent.Packages) map[string]audiocontent.Packages

func BuildFilterPipeline(cmd *cobra.Command) Filter {
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
