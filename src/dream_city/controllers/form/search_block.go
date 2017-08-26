package form

type FSearchBlock struct {
	BlockName     string
	CityName      string
	UserName      string
}

func (f *FSearchBlock) Valid() bool {
	return f.BlockName != "" || f.CityName != "" || f.UserName != ""
}

func (f *FSearchBlock) Labels() map[string]string {
	return defaultColumnNames(f)
}

