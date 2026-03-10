package models

type Tutor struct {
	UserID      string
	Name        string
	HourlyRate  int
	Description string
	PhotoPath   string
	SubjectIDs  []int64
}
