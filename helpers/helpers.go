package helpers

import (
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
