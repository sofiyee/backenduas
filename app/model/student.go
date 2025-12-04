package model

import "time"

type Student struct {
	UserID        string `json:"user_id"`
	StudentID     string `json:"student_id"`
	ProgramStudy  string `json:"program_study"`
	AcademicYear  string `json:"academic_year"`
	AdvisorID     *string `json:"advisor_id"`
}

type CreateStudentRequest struct {
	UserID        string  `json:"user_id"`
	StudentID     string  `json:"student_id"`
	ProgramStudy  string  `json:"program_study"`
	AcademicYear  string  `json:"academic_year"`
	AdvisorID     *string `json:"advisor_id"`
}

type StudentDetail struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	StudentID     string    `json:"student_id"`
	ProgramStudy  string    `json:"program_study"`
	AcademicYear  string    `json:"academic_year"`
	AdvisorID     *string   `json:"advisor_id"`
	AdvisorName   *string   `json:"advisor_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type StudentAchievement struct {
    ID     string        `json:"id"`
    Status string        `json:"status"`
    Detail AchievementMongo `json:"detail"`
}
