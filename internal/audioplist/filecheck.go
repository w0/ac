package audioplist

import "fmt"

type FileCheck []string

func (fc *FileCheck) UnmarshalPlist(unmarshal func(any) error) error {
	var fcString string

	err := unmarshal(&fcString)

	if err == nil {
		*fc = []string{fcString}
		return nil
	}

	var fcArray []string
	err = unmarshal(&fcArray)

	if err == nil {
		*fc = fcArray
		return nil
	}

	return fmt.Errorf("failed to unmarshal as string or []string: %w", err)
}
