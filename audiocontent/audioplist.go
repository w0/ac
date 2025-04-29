package audiocontent

type Packages struct {
	ContainsAppleLoops                  bool         `plist:"ContainsAppleLoops"`
	ContainsGarageBandLegacyInstruments bool         `plist:"ContainsGarageBandLegacyInstruments"`
	DownloadName                        DownloadName `plist:"DownloadName"`
	DownloadSize                        ContentSize  `plist:"DownloadSize"`
	FileCheck                           FileCheck    `plist:"FileCheck"`
	InstalledSize                       ContentSize  `plist:"InstalledSize"`
	IsMandatory                         bool         `plist:"IsMandatory"`
	PackageID                           string       `plist:"PackageID"`
}

type Content struct {
	Packages map[string]Packages `plist:"Packages"`
}

func (ac Content) GetPackageByName(names []string) []Packages {
	var pkgs []Packages

	for _, name := range names {
		pkgs = append(pkgs, ac.Packages[name])
	}

	return pkgs

}
