package repository

import (
	"context"
	"backenduas/app/model"
	"backenduas/database"
)

type LecturerRepository struct {}

func NewLecturerRepository() *LecturerRepository {
	return &LecturerRepository{}
}

// GET /lecturers
func (r *LecturerRepository) GetAll(ctx context.Context) ([]model.LecturerDetail, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT l.id, l.user_id, u.username, u.email, u.full_name,
		       l.lecturer_id, l.department, l.created_at
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY l.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.LecturerDetail

	for rows.Next() {
		var d model.LecturerDetail
		if err := rows.Scan(
			&d.ID, &d.UserID, &d.Username, &d.Email, &d.FullName,
			&d.LecturerID, &d.Department, &d.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

// GET /lecturers/:id/advisees
func (r *LecturerRepository) GetAdvisees(ctx context.Context, lecturerID string) ([]model.Advisee, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT s.student_id, u.full_name, s.program_study, s.academic_year
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Advisee

	for rows.Next() {
		var a model.Advisee
		if err := rows.Scan(&a.StudentID, &a.StudentName, &a.ProgramStudy, &a.AcademicYear); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}
