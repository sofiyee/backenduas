package model

type Advisee struct {
	StudentID     string `json:"student_id"`
	StudentName   string `json:"student_name"`
	ProgramStudy  string `json:"program_study"`
	AcademicYear  string `json:"academic_year"`
}
