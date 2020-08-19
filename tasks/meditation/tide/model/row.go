package model

type MusicRow struct {
	Id             string
	Name           string
	DemoLink       string
	CoverLink      string
	Duration       int
	Description    string
	SubTitle       string
	PrimaryColor   string
	SecondaryColor string
	UpdatedAt      int64
	CreatedAt      int64
	TagIds         string
	Loved          int
	Hash           string
}

type TagRow struct {
	TagId   string
	SortKey int
	Key     string
	Type    string
	Name    string
}

type MedRow struct {
	Type         int
	MedId        string
	SectionId    string
	ResId        string
	TagIds       string
	Name         string
	Description  string
	PrimaryColor string
	CreatedAt    int64
	UpdatedAt    int64
	SortKey      int
	DemoLink     string
	CoverLink    string
	Duration     int
	Hash         string
	ResLink      string
}
