package model

type TopStudent struct {
    StudentID string `json:"student_id"`
    Count     int    `json:"count"`
}

type StudentAchievementStat struct {
    PerType      map[string]int      `json:"per_type"`
    PerPeriod    map[string]int      `json:"per_period"`
    Competition  map[string]int      `json:"competition_level"`
    TopStudents  []TopStudentDetail  `json:"top_students"`  
}

type TopStudentDetail struct {
    StudentID     string `json:"student_id"`
    FullName      string `json:"full_name"`
    ProgramStudy  string `json:"program_study"`
    AcademicYear  string `json:"academic_year"`
    Count         int    `json:"count"`
}
