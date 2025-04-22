package audiocontent

type ContentSize int

func (cs *ContentSize) UnmarshalPlist(unmarshal func(any) error) error {

	var csInt int
	err := unmarshal(&csInt)

	if err == nil {
		*cs = ContentSize(csInt)
		return nil
	}

	var csFloat float64
	if err := unmarshal(&csFloat); err != nil {
		return err
	}

	*cs = ContentSize(int(csFloat))

	return nil
}
