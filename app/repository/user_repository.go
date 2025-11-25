package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"backenduas/app/model"
	"backenduas/database"

	"github.com/google/uuid"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT u.id, u.username, u.email, u.full_name,
			   u.role_id, r.name AS role_name,
			   u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
			&u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT u.id, u.username, u.email, u.full_name,
			   u.role_id, r.name AS role_name,
			   u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
	`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName,
		&u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {

	// INSERT USER
	_, err := database.DB.Exec(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, true, NOW(), NOW())
	`, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID)
	if err != nil {
		return err
	}

	// GET USER ID
	err = database.DB.QueryRow(ctx, `SELECT id FROM users WHERE username=$1`, u.Username).Scan(&u.ID)
	if err != nil {
		return err
	}

	// CEK ROLE
	roleName, err := r.GetRoleNameByID(ctx, u.RoleID)
	if err != nil {
		return err
	}
	roleName = strings.ToLower(roleName)

	// ===========================================
	// AUTO INSERT DOSEN WALI
	// ===========================================
	if roleName == "dosen wali" {

		newLecturerID, _ := r.GenerateNextLecturerID(ctx)
		lecturerUUID := uuid.New().String()

		_, err := database.DB.Exec(ctx, `
			INSERT INTO lecturers (id, user_id, lecturer_id, department, created_at)
			VALUES ($1, $2, $3, $4, NOW())
		`, lecturerUUID, u.ID, newLecturerID, "Teknik Informatika")

		if err != nil {
			return err
		}
	}

	// ===========================================
	// AUTO INSERT MAHASISWA
	// ===========================================
if roleName == "mahasiswa" {

	nextStudentID, _ := r.GenerateNextStudentID(ctx)

	advisorID, _ := r.GetRandomAdvisorID(ctx)

	studentUUID := uuid.New().String()

	// FIX: jika advisor kosong beri NULL
	var advisor interface{}
	if advisorID == "" {
		advisor = nil
	} else {
		advisor = advisorID
	}

	_, err := database.DB.Exec(ctx, `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, studentUUID, u.ID, nextStudentID, "Teknik Informatika", "2023", advisor)

	if err != nil {
		return err
	}
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE users SET username=$1, email=$2, full_name=$3, updated_at=NOW()
		WHERE id = $4
	`, u.Username, u.Email, u.FullName, u.ID)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := database.DB.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID string, roleID string) error {
	_, err := database.DB.Exec(ctx, `
		UPDATE users SET role_id=$1, updated_at=NOW()
		WHERE id = $2
	`, roleID, userID)
	return err
}

// =======================
// GET ROLE NAME
// =======================
func (r *UserRepository) GetRoleNameByID(ctx context.Context, roleID string) (string, error) {
	var roleName string
	err := database.DB.QueryRow(ctx, `SELECT name FROM roles WHERE id=$1`, roleID).Scan(&roleName)
	if err != nil {
		return "", err
	}
	return roleName, nil
}

// =======================
// NEXT STUDENT ID
// =======================
func (r *UserRepository) GenerateNextStudentID(ctx context.Context) (string, error) {
	var last string

	err := database.DB.QueryRow(ctx, `
		SELECT student_id FROM students ORDER BY created_at DESC LIMIT 1
	`).Scan(&last)

	if err != nil {
		return "S0001", nil 
	}

	num, _ := strconv.Atoi(last[1:])
	next := fmt.Sprintf("S%04d", num+1)

	return next, nil
}

// =======================
// NEXT LECTURER ID
// =======================
func (r *UserRepository) GenerateNextLecturerID(ctx context.Context) (string, error) {
	var last string

	err := database.DB.QueryRow(ctx, `
		SELECT lecturer_id FROM lecturers ORDER BY created_at DESC LIMIT 1
	`).Scan(&last)

	if err != nil {
		return "L0001", nil
	}

	num, _ := strconv.Atoi(last[1:])
	next := fmt.Sprintf("L%04d", num+1)

	return next, nil
}

// =======================
// RANDOM ADVISOR
// =======================
func (r *UserRepository) GetRandomAdvisorID(ctx context.Context) (string, error) {
	var id string

	err := database.DB.QueryRow(ctx, `
		SELECT id FROM lecturers ORDER BY RANDOM() LIMIT 1
	`).Scan(&id)

	if err != nil {
		return "", nil // boleh kosong
	}

	return id, nil
}
