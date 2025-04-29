package audiocontent

import (
	"fmt"
	"strings"
)

type DownloadName string

func (dn *DownloadName) UnmarshalPlist(unmarshal func(any) error) error {
	acBaseUrl := "https://audiocontentdownload.apple.com"
	acFilePath := "/lp10_ms3_content_2016/"

	var dnString string
	err := unmarshal(&dnString)

	if err != nil {
		return err
	}

	if strings.Contains(dnString, "../lp10_ms3_content_2013/") {
		striped := strings.TrimPrefix(dnString, "..")
		*dn = DownloadName(fmt.Sprintf("%s%s", acBaseUrl, striped))
	} else {
		*dn = DownloadName(fmt.Sprintf("%s%s%s", acBaseUrl, acFilePath, dnString))
	}

	return nil

}
