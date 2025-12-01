package repository

import (
	"context"
	"backenduas/app/model"
	"backenduas/database"
)

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{}
}

// ===============================
// GET ALL STUDENTS
// ===============================
func (r *StudentRepository) GetAll(ctx context.Context) ([]model.StudentDetail, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT 
			s.id,
			s.user_id,
			u.username,
			u.email,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year,
			s.advisor_id,
			l.full_name AS advisor_name,
			s.created_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers lc ON lc.id = s.advisor_id
		LEFT JOIN users l ON l.id = lc.user_id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.StudentDetail

	for rows.Next() {
		var s model.StudentDetail
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Username,
			&s.Email,
			&s.FullName,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.AdvisorName,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

// ===============================
// GET STUDENT BY ID
// ===============================
func (r *StudentRepository) GetByID(ctx context.Context, id string) (*model.StudentDetail, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT 
			s.id,
			s.user_id,
			u.username,
			u.email,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year,
			s.advisor_id,
			l.full_name AS advisor_name,
			s.created_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers lc ON lc.id = s.advisor_id
		LEFT JOIN users l ON l.id = lc.user_id
		WHERE s.id = $1
	`, id)

	var s model.StudentDetail
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.Username,
		&s.Email,
		&s.FullName,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.AdvisorName,
		&s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// ===============================
// CREATE STUDENT (ADMIN)
// ===============================
func (r *StudentRepository) Create(ctx context.Context, s *model.Student) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
	`,
		s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID,
	)

	return err
}

// ===============================
// UPDATE ADVISOR
// ===============================
func (r *StudentRepository) UpdateAdvisor(ctx context.Context, studentID string, advisorID *string) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`, advisorID, studentID)

	return err
}

// ===============================
// GET STUDENT ID BY USER ID (Mahasiswa Login)
// ===============================
func (r *StudentRepository) GetStudentIDByUserID(ctx context.Context, userID string) (string, error) {
	row := database.DB.QueryRow(ctx, `
        SELECT id 
        FROM students
        WHERE user_id = $1
    `, userID)

	var studentID string
	err := row.Scan(&studentID)
	return studentID, err
}

// ===============================
// GET LECTURER ID BY USER ID (Dosen Login)
// ===============================
func (r *StudentRepository) GetLecturerIDByUserID(ctx context.Context, userID string) (string, error) {
	row := database.DB.QueryRow(ctx, `
        SELECT id 
        FROM lecturers
        WHERE user_id = $1
    `, userID)

	var lecturerID string
	err := row.Scan(&lecturerID)
	return lecturerID, err
}

// ===============================
// GET STUDENTS BY ADVISOR (Mahasiswa Bimbingan Dosen)
// ===============================
func (r *StudentRepository) GetStudentsByAdvisor(ctx context.Context, advisorID string) ([]string, error) {
	rows, err := database.DB.Query(ctx, `
        SELECT id
        FROM students
        WHERE advisor_id = $1
    `, advisorID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	studentIDs := []string{}

	for rows.Next() {
		var id string
		rows.Scan(&id)
		studentIDs = append(studentIDs, id)
	}

	return studentIDs, nil
}

func (r *StudentRepository) IsStudentUnderAdvisor(ctx context.Context, studentID string, advisorID string) (bool, error) {
    row := database.DB.QueryRow(ctx, `
        SELECT COUNT(*) 
        FROM students 
        WHERE id = $1 AND advisor_id = $2
    `, studentID, advisorID)

    var count int
    err := row.Scan(&count)
    return count > 0, err
}
