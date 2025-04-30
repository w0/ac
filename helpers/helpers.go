package helpers

import (
	"fmt"
	"os"

	"github.com/micromdm/plist"
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

func PrettyPrint(names []string, ac audiocontent.Content) {
	for _, name := range names {
		fmt.Printf("\n-- %s --\n", name)

		if pkg, ok := ac.Packages[name]; ok {
			fmt.Printf("\t- ContainsAppleLoops: %t\n", pkg.ContainsAppleLoops)
			fmt.Printf("\t- ContainsGarageBandLegacyInstruments: %t\n", pkg.ContainsGarageBandLegacyInstruments)
			fmt.Printf("\t- DownloadName: %s\n", pkg.DownloadName)
			fmt.Printf("\t- DownloadSize: %0.2f Mb\n", ConvertToMB(pkg.DownloadSize))
			fmt.Printf("\t- FileCheck:\n")

			for _, fc := range pkg.FileCheck {
				fmt.Printf("\t\t- %s\n", fc)
			}

			fmt.Printf("\t- InstalledSize: %0.2f Mb\n", ConvertToMB(pkg.InstalledSize))
			fmt.Printf("\t- IsMandatory: %t\n", pkg.IsMandatory)
			fmt.Printf("\t- PackageID: %s\n", pkg.PackageID)
		}

	}
}

func ConvertToMB(n audiocontent.ContentSize) float64 {
	return float64(n) / (1 << 20)
}
