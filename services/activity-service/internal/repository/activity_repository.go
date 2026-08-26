package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"activity-service/internal/model"
)

type ActivityRepository interface {
	Create(activity *model.DailyActivity) (*model.DailyActivity, error)
	GetByFilter(userID int32, userNIP, status string, limit int32) ([]*model.DailyActivity, error)
	GetByID(id int32) (*model.DailyActivity, error)
	UpdateStatus(id int32, newStatus string) error
	Delete(id int32, userID int32) error
}

type activityRepository struct {
	db        *sql.DB
	memStore  map[int32]*model.DailyActivity
	autoIncID int32
}

func NewActivityRepository(db *sql.DB) ActivityRepository {
	repo := &activityRepository{
		db:        db,
		memStore:  make(map[int32]*model.DailyActivity),
		autoIncID: 1,
	}

	if db == nil {
		repo.initSeedData()
	}

	return repo
}

func (r *activityRepository) initSeedData() {
	seed := []*model.DailyActivity{
		{
			ID:           1,
			UserID:       1,
			UserNIP:      "IN-294",
			UserName:     "Sarah Jenkins",
			Title:        "Penginputan Data Laporan Harian Posko 112",
			Description:  "Melakukan verifikasi dan penginputan 15 laporan darurat masyarakat ke dalam database sistem LOPI-Q.",
			ActivityDate: "2026-08-26",
			PhotoURL:     "",
			Status:       "APPROVED",
			CreatedAt:    time.Now().Add(-2 * time.Hour),
			UpdatedAt:    time.Now().Add(-2 * time.Hour),
		},
		{
			ID:           2,
			UserID:       2,
			UserNIP:      "IN-291",
			UserName:     "Marcus Doe",
			Title:        "Pemeliharaan Perangkat Jaringan Siaga",
			Description:  "Pemeriksaan rutin kabel LAN dan router komunikasi darurat di ruang operator Posko NTPD 112.",
			ActivityDate: "2026-08-26",
			PhotoURL:     "",
			Status:       "PENDING",
			CreatedAt:    time.Now().Add(-1 * time.Hour),
			UpdatedAt:    time.Now().Add(-1 * time.Hour),
		},
	}

	for _, item := range seed {
		r.memStore[item.ID] = item
		r.autoIncID = item.ID + 1
	}
}

func (r *activityRepository) Create(act *model.DailyActivity) (*model.DailyActivity, error) {
	act.CreatedAt = time.Now()
	act.UpdatedAt = time.Now()
	if act.Status == "" {
		act.Status = "PENDING"
	}

	if r.db != nil {
		query := `
			INSERT INTO tb_lopiq_daily_activities 
			(user_id, user_nip, user_name, title, description, activity_date, photo_url, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`
		err := r.db.QueryRow(
			query,
			act.UserID, act.UserNIP, act.UserName, act.Title, act.Description,
			act.ActivityDate, act.PhotoURL, act.Status, act.CreatedAt, act.UpdatedAt,
		).Scan(&act.ID)

		if err != nil {
			log.Printf("Error inserting activity into DB: %v", err)
			return nil, err
		}
		return act, nil
	}

	// Memory fallback
	act.ID = r.autoIncID
	r.autoIncID++
	r.memStore[act.ID] = act
	return act, nil
}

func (r *activityRepository) GetByFilter(userID int32, userNIP, status string, limit int32) ([]*model.DailyActivity, error) {
	if limit <= 0 {
		limit = 50
	}

	if r.db != nil {
		query := `SELECT id, user_id, user_nip, user_name, title, description, activity_date, photo_url, status, created_at, updated_at FROM tb_lopiq_daily_activities WHERE 1=1`
		args := []interface{}{}
		argIdx := 1

		if userID > 0 {
			query += fmt.Sprintf(" AND user_id = $%d", argIdx)
			args = append(args, userID)
			argIdx++
		}
		if userNIP != "" {
			query += fmt.Sprintf(" AND user_nip = $%d", argIdx)
			args = append(args, userNIP)
			argIdx++
		}
		if status != "" {
			query += fmt.Sprintf(" AND status = $%d", argIdx)
			args = append(args, status)
			argIdx++
		}

		query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d", argIdx)
		args = append(args, limit)

		rows, err := r.db.Query(query, args...)
		if err != nil {
			log.Printf("Error querying activities: %v", err)
			return nil, err
		}
		defer rows.Close()

		var results []*model.DailyActivity
		for rows.Next() {
			item := &model.DailyActivity{}
			if err := rows.Scan(
				&item.ID, &item.UserID, &item.UserNIP, &item.UserName,
				&item.Title, &item.Description, &item.ActivityDate,
				&item.PhotoURL, &item.Status, &item.CreatedAt, &item.UpdatedAt,
			); err != nil {
				continue
			}
			results = append(results, item)
		}
		return results, nil
	}

	// Memory fallback
	var results []*model.DailyActivity
	for _, item := range r.memStore {
		if userID > 0 && item.UserID != userID {
			continue
		}
		if userNIP != "" && item.UserNIP != userNIP {
			continue
		}
		if status != "" && item.Status != status {
			continue
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *activityRepository) GetByID(id int32) (*model.DailyActivity, error) {
	if r.db != nil {
		query := `SELECT id, user_id, user_nip, user_name, title, description, activity_date, photo_url, status, created_at, updated_at FROM tb_lopiq_daily_activities WHERE id = $1`
		item := &model.DailyActivity{}
		err := r.db.QueryRow(query, id).Scan(
			&item.ID, &item.UserID, &item.UserNIP, &item.UserName,
			&item.Title, &item.Description, &item.ActivityDate,
			&item.PhotoURL, &item.Status, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return item, nil
	}

	item, exists := r.memStore[id]
	if !exists {
		return nil, fmt.Errorf("activity not found with id %d", id)
	}
	return item, nil
}

func (r *activityRepository) UpdateStatus(id int32, newStatus string) error {
	if r.db != nil {
		query := `UPDATE tb_lopiq_daily_activities SET status = $1, updated_at = $2 WHERE id = $3`
		_, err := r.db.Exec(query, newStatus, time.Now(), id)
		return err
	}

	item, exists := r.memStore[id]
	if !exists {
		return fmt.Errorf("activity not found with id %d", id)
	}
	item.Status = newStatus
	item.UpdatedAt = time.Now()
	return nil
}

func (r *activityRepository) Delete(id int32, userID int32) error {
	if r.db != nil {
		query := `DELETE FROM tb_lopiq_daily_activities WHERE id = $1`
		if userID > 0 {
			query += fmt.Sprintf(" AND user_id = %d", userID)
		}
		_, err := r.db.Exec(query, id)
		return err
	}

	delete(r.memStore, id)
	return nil
}
