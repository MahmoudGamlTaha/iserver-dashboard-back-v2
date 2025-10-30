package repositories

import (
	"database/sql"
	"enterprise-architect-api/models"
	"fmt"
	"strings"
	"time"
)

// ProfileRepository handles database operations for profiles
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository creates a new ProfileRepository
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create creates a new profile in the database
func (r *ProfileRepository) Create(req models.CreateProfileRequest) (*models.Profile, error) {
	now := time.Now()

	query := `
		INSERT INTO Profile (
			ProfileName, ProfileDescription, PortalStartPageId, DateCreated, CreatedBy, DateModified, ModifiedBy
		) OUTPUT INSERTED.ProfileID VALUES (
			@p1, @p2, @p3, @p4, @p5, @p6, @p7
		)
	`

	var profileID int
	err := r.db.QueryRow(query,
		req.ProfileName, req.ProfileDescription, req.PortalStartPageId, now, req.CreatedBy, now, req.CreatedBy,
	).Scan(&profileID)

	if err != nil {
		return nil, fmt.Errorf("error creating profile: %w", err)
	}

	return r.GetByID(profileID)
}

// GetByID retrieves a profile by its ID
func (r *ProfileRepository) GetByID(id int) (*models.Profile, error) {
	query := `
		SELECT ProfileID, ProfileName, ProfileDescription, PortalStartPageId, DateCreated, CreatedBy, DateModified, ModifiedBy
		FROM Profile
		WHERE ProfileID = @p1
	`

	profile := &models.Profile{}
	err := r.db.QueryRow(query, id).Scan(
		&profile.ProfileID, &profile.ProfileName, &profile.ProfileDescription, &profile.PortalStartPageId,
		&profile.DateCreated, &profile.CreatedBy, &profile.DateModified, &profile.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving profile: %w", err)
	}

	return profile, nil
}

// GetAll retrieves all profiles with pagination
func (r *ProfileRepository) GetAll(page, pageSize int) ([]models.Profile, int, error) {
	offset := (page - 1) * pageSize

	// Get total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM Profile`
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting profiles: %w", err)
	}

	// Get paginated results
	query := `
		SELECT ProfileID, ProfileName, ProfileDescription, PortalStartPageId, DateCreated, CreatedBy, DateModified, ModifiedBy
		FROM Profile
		ORDER BY DateCreated DESC
		OFFSET @p1 ROWS FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := r.db.Query(query, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile
		err := rows.Scan(
			&profile.ProfileID, &profile.ProfileName, &profile.ProfileDescription, &profile.PortalStartPageId,
			&profile.DateCreated, &profile.CreatedBy, &profile.DateModified, &profile.ModifiedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	return profiles, totalCount, nil
}

// Update updates an existing profile
func (r *ProfileRepository) Update(id int, req models.UpdateProfileRequest) (*models.Profile, error) {
	// Build dynamic update query
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.ProfileName != nil {
		setClauses = append(setClauses, fmt.Sprintf("ProfileName = @p%d", argIndex))
		args = append(args, *req.ProfileName)
		argIndex++
	}
	if req.ProfileDescription != nil {
		setClauses = append(setClauses, fmt.Sprintf("ProfileDescription = @p%d", argIndex))
		args = append(args, *req.ProfileDescription)
		argIndex++
	}
	if req.PortalStartPageId != nil {
		setClauses = append(setClauses, fmt.Sprintf("PortalStartPageId = @p%d", argIndex))
		args = append(args, *req.PortalStartPageId)
		argIndex++
	}

	// Always update DateModified and ModifiedBy
	setClauses = append(setClauses, fmt.Sprintf("DateModified = @p%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	setClauses = append(setClauses, fmt.Sprintf("ModifiedBy = @p%d", argIndex))
	args = append(args, req.ModifiedBy)
	argIndex++

	if len(setClauses) == 2 { // Only DateModified and ModifiedBy
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE Profile SET %s WHERE ProfileID = @p%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating profile: %w", err)
	}

	return r.GetByID(id)
}

// Delete deletes a profile by its ID
func (r *ProfileRepository) Delete(id int) error {
	query := `DELETE FROM Profile WHERE ProfileID = @p1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}

