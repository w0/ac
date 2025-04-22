package audiocontent

import (
	"testing"

	"github.com/micromdm/plist"
)

const testSizeData = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Packages</key>
	<dict>
		<key>Integer</key>
		<dict>
			<key>DownloadSize</key>
			<integer>506768673</integer>
		</dict>
		<key>Float</key>
		<dict>
			<key>DownloadSize</key>
			<real>201060739</real>
		</dict>
	</dict>
</dict>
</plist>
`

func TestAudioSize(t *testing.T) {

	type packageInfo struct {
		DownloadSize ContentSize `plist:"DownloadSize"`
	}

	var packages struct {
		Packages map[string]packageInfo `plist:"Packages"`
	}

	err := plist.Unmarshal([]byte(testSizeData), &packages)

	if err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}

	if p := packages.Packages["Integer"].DownloadSize; p != 506768673 {
		t.Fatalf("DownloadSize is incorrect. Got: %d", p)
	}

}
