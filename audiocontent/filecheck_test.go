package audiocontent

import (
	"testing"

	"github.com/micromdm/plist"
)

const testFileData = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Packages</key>
	<dict>
		<key>String</key>
		<dict>
			<key>FileCheck</key>
			<string>/Library/Application Support/GarageBand/Instrument Library/Track Settings/Software/Woodwinds/Tenor Sax.cst</string>
		</dict>
		<key>Array</key>
		<dict>
			<key>FileCheck</key>
			<array>
				<string>/Library/Application Support/GarageBand/Apple Loops/Apple Loops for GarageBand Jam Pack/12 String Dream 01.aif</string>
				<string>/Library/Audio/Apple Loops/Apple/Jam Pack 1/12 String Dream 01.aif</string>
				<string>/Library/Audio/Apple Loops/Apple/Jam Pack 1/12 String Dream 01.caf</string>
			</array>
		</dict>
	</dict>
</dict>
</plist>
`

func TestFileCheck(t *testing.T) {
	type fcTest struct {
		Packages map[string]struct {
			FileCheck FileCheck `plist:"FileCheck"`
		} `plist:"Packages"`
	}

	var test fcTest

	err := plist.Unmarshal([]byte(testFileData), &test)

	if err != nil {
		t.Fail()
	}

}
