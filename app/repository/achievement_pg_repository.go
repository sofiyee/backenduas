package repository

import (
	"context"
	"backenduas/app/model"
	"backenduas/database"
)

type AchievementPGRepository struct{}

func NewAchievementPGRepository() *AchievementPGRepository {
	return &AchievementPGRepository{}
}

func (r *AchievementPGRepository) CreateReference(ref model.AchievementReference) error {
	_, err := database.DB.Exec(context.Background(), `
		INSERT INTO achievement_references
			(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,NOW(),NOW())
	`,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
	)
	return err
}

func (r *AchievementPGRepository) UpdateStatus(id string, status string) error {
	ctx := context.Background()

	// jika submitted → update submitted_at dulu
	if status == "submitted" {
		_, err := database.DB.Exec(ctx, `
			UPDATE achievement_references
			SET submitted_at = NOW()
			WHERE id = $1
		`, id)
		if err != nil {
			return err
		}
	}

	// update status biasa (harus cast ENUM)
	_, err := database.DB.Exec(ctx, `
		UPDATE achievement_references
		SET status = $1::achievement_status,
		    updated_at = NOW()
		WHERE id = $2
	`, status, id)

	return err
}


func (r *AchievementPGRepository) Verify(id, verifier string) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE achievement_references
		SET status='verified', verified_by=$1, verified_at=NOW()
		WHERE id=$2
	`, verifier, id)
	return err
}

func (r *AchievementPGRepository) Reject(id, note string) error {
	_, err := database.DB.Exec(context.Background(), `
		UPDATE achievement_references
		SET status='rejected', rejection_note=$1, updated_at=NOW()
		WHERE id=$2
	`, note, id)
	return err
}

func (r *AchievementPGRepository) GetByID(ctx context.Context, id string) (model.AchievementReference, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`, id)

	var ref model.AchievementReference
	err := row.Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	return ref, err
}

// GET by student (mahasiswa)
func (r *AchievementPGRepository) GetByStudentID(ctx context.Context, studentID string) ([]model.AchievementReference, error) {
	rows, err := database.DB.Query(ctx, `
        SELECT id, student_id, mongo_achievement_id, status
        FROM achievement_references
        WHERE student_id = $1
    `, studentID)

	if err != nil {
		return nil, err
	}

	list := []model.AchievementReference{}

	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status)
		list = append(list, ref)
	}

	return list, nil
}

// GET all (admin)
func (r *AchievementPGRepository) GetAll(ctx context.Context) ([]model.AchievementReference, error) {
	rows, err := database.DB.Query(ctx, `
        SELECT id, student_id, mongo_achievement_id, status
        FROM achievement_references
    `)

	if err != nil {
		return nil, err
	}

	list := []model.AchievementReference{}

	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status)
		list = append(list, ref)
	}

	return list, nil
}

// ⭐ UNTUK DOSEN WALI: Ambil achievement dari list student_ids
func (r *AchievementPGRepository) GetByStudentIDs(ctx context.Context, studentIDs []string) ([]model.AchievementReference, error) {

	rows, err := database.DB.Query(ctx, `
        SELECT id, student_id, mongo_achievement_id, status
        FROM achievement_references
        WHERE student_id = ANY($1)
    `, studentIDs)

	if err != nil {
		return nil, err
	}

	list := []model.AchievementReference{}

	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status)
		list = append(list, ref)
	}

	return list, nil
}

